// Package storage 提供本地文件存储。
//
// 存储只有本地磁盘一种实现，目录布局为 <root>/<bucket>/<key>，
// 与对外 URL /<bucket>/<key> 一一对应。路径与桶名固定，无需运行时配置，
// 数据库中已保存的 key 也不需要迁移。
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codfrm/cago/configs"
)

// 存储位置固定，随容器的 NFS 卷挂载到 /data
const (
	defaultRoot   = "/data"
	defaultBucket = "dsp2b"
)

// ErrNotFound 对象不存在
var ErrNotFound = errors.New("object not found")

// RespError 带状态码的存储错误，与对象存储的错误语义保持一致
// 调用方用 StatusCode() == 404 判断对象不存在
type RespError interface {
	error
	StatusCode() int
}

// notFoundError 满足调用方对 RespError 的断言（StatusCode == 404）
type notFoundError struct{}

func (notFoundError) Error() string   { return "object not found" }
func (notFoundError) StatusCode() int { return 404 }

// Is 让 errors.Is(err, ErrNotFound) 对 notFoundError 生效
func (notFoundError) Is(target error) bool { return target == ErrNotFound }

var defaultStore *Store

// Store 本地文件存储
type Store struct {
	root   string
	bucket string
}

// Default 获取默认存储实例
func Default() *Store {
	return defaultStore
}

// Init 初始化默认存储，并确保桶目录存在
// 签名保持与 cago 组件一致，存储位置固定，忽略传入的上下文与配置
func Init(ctx context.Context, config *configs.Config) error {
	dir := filepath.Join(defaultRoot, defaultBucket)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create storage dir %s: %w", dir, err)
	}
	defaultStore = &Store{root: defaultRoot, bucket: defaultBucket}
	return nil
}

// BucketDir 返回桶目录
func (s *Store) BucketDir() string {
	return filepath.Join(s.root, s.bucket)
}

// resolve 把 object key 解析成磁盘路径，并阻止路径穿越
func (s *Store) resolve(key string) (string, error) {
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return "", fmt.Errorf("empty object key")
	}
	base := s.BucketDir()
	full := filepath.Join(base, filepath.FromSlash(key))
	if full != base && !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid object key %q", key)
	}
	return full, nil
}

// PutObject 写入对象，先写临时文件再 rename，避免读到半截内容
func (s *Store) PutObject(ctx context.Context, key string, r io.Reader, size int64) error {
	file, err := s.resolve(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	tmp := file + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, file)
}

// Object 打开的本地对象
type Object struct {
	*os.File
}

// Stat 返回对象信息
func (o *Object) Stat() (ObjectInfo, error) {
	st, err := o.File.Stat()
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{
		Size:         st.Size(),
		LastModified: st.ModTime(),
	}, nil
}

// ObjectInfo 对象元信息
type ObjectInfo struct {
	Key          string
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

// GetObject 读取对象
func (s *Store) GetObject(ctx context.Context, key string) (*Object, error) {
	file, err := s.resolve(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, notFoundError{}
		}
		return nil, err
	}
	return &Object{File: f}, nil
}

// Exists 判断对象是否存在
func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	file, err := s.resolve(key)
	if err != nil {
		return false, err
	}
	st, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !st.IsDir(), nil
}
