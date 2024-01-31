package command

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func cloneCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		logger.Default().Error("参数错误")
		return errors.New("参数错误")
	}
	id, err := primitive.ObjectIDFromHex(args[0])
	if err != nil {
		logger.Default().Error("参数错误")
		return errors.New("参数错误")
	}
	// 存在dspb.json
	if _, err := os.Stat("dspb.json"); err == nil {
		logger.Default().Error("仓库已存在")
		return errors.New("仓库已存在")
	} else if !os.IsNotExist(err) {
		logger.Default().Error("判断文件是否存在失败", zap.Error(err))
		return err
	}
	if err := cloneRepo(id); err != nil {
		return err
	}
	return nil
}

func cloneRepo(id primitive.ObjectID) error {
	// 写dspb.json
	ctx := context.Background()
	logger.Ctx(ctx).Info("加载蓝图集中...", zap.String("id", id.Hex()))
	collection, err := CollectionDetail(ctx, id, 1)
	if err != nil {
		return err
	}
	var data []byte
	// 下载蓝图集zip
	logger.Ctx(ctx).Info("下载蓝图集zip中...", zap.String("url", collection.Collection.DownloadFile))
	downloadResp, err := http.Get(collection.Collection.DownloadFile)
	if err != nil {
		return err
	}
	defer downloadResp.Body.Close()
	if downloadResp.StatusCode != http.StatusOK {
		logger.Ctx(ctx).Info("请求失败", zap.String("body", string(data)))
		return errors.New("请求失败")
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, downloadResp.Body); err != nil {
		return err
	}
	data = buf.Bytes()
	buf.Reset()
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	// 解压zip
	logger.Ctx(ctx).Info("解压zip中...")
	var hashMap = make(map[string]*Blueprint)
	for _, file := range zipReader.File {
		err := func(file *zip.File) error {
			if file.FileInfo().IsDir() {
				// 本地也创建目录
				return nil
			}
			logger.Ctx(ctx).Info("解压文件", zap.String("name", file.Name))
			f, err := file.Open()
			if err != nil {
				return err
			}
			defer f.Close()
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, f); err != nil {
				return err
			}
			// 计算一次hash
			hash, err := Hash(buf.Bytes())
			if err != nil {
				return err
			}
			// 记录hash
			hashMap[file.Name] = &Blueprint{
				Hash:            hash,
				Updatetime:      file.Modified.Unix(),
				LocalUpdatetime: file.Modified.Unix(),
			}
			// 判断本地文件是否存在，如果存在则跳过，不存在则解压
			stat, err := os.Stat(file.Name)
			if err != nil {
				if os.IsNotExist(err) {
					// 不存在
					// 创建目录
					if dir := filepath.Dir(file.Name); dir != "" {
						if err := os.MkdirAll(dir, os.ModePerm); err != nil {
							return err
						}
					}
					newfile, err := os.Create(file.Name)
					if err != nil {
						return err
					}
					defer newfile.Close()
					if _, err := newfile.Write(buf.Bytes()); err != nil {
						return err
					}
					return nil
				}
				logger.Ctx(ctx).Error("判断文件是否存在失败", zap.Error(err),
					zap.String("name", file.Name))
				return err
			}
			hashMap[file.Name].LocalUpdatetime = stat.ModTime().Unix()
			// 存在跳过
			return nil
		}(file)
		if err != nil {
			logger.Ctx(ctx).Error("解压文件失败", zap.Error(err), zap.String("name", file.Name))
			return err
		}
	}
	// 构建dspb.json
	logger.Ctx(ctx).Info("构建dspb.json中...")
	repo := &Repository{
		ID:          id,
		Title:       collection.Collection.Title,
		Description: collection.Collection.Description,
		Version:     "1.0.0",
		Blueprint:   make([]*Blueprint, 0),
		Repository:  make([]*Repository, 0),
	}

	if err := bfs(ctx, "", repo, id, hashMap); err != nil {
		logger.Ctx(ctx).Info("构建dspb.json失败", zap.Error(err))
		return err
	}
	// 写入文件
	logger.Ctx(ctx).Info("写入dspb.json中...")
	if err := SaveRepository(repo); err != nil {
		return err
	}
	return nil
}

// 遍历蓝图集，生成dspb.json
func bfs(ctx context.Context, path string, repo *Repository, id primitive.ObjectID, hashMap map[string]*Blueprint) error {
	info, items, err := ReadAllBlueprint(ctx, id)
	if err != nil {
		return err
	}
	for _, item := range items {
		filename := filepath.Join(path, filepath.Clean(item.Title)+".txt")
		hash, ok := hashMap[filename]
		if !ok {
			logger.Ctx(ctx).Error("未找到hash", zap.String("title", item.Title))
			return errors.New("未找到hash")
		}
		repo.Blueprint = append(repo.Blueprint, &Blueprint{
			ID:              id,
			Title:           item.Title,
			Description:     item.Description,
			Hash:            hash.Hash,
			Updatetime:      hash.Updatetime,
			LocalUpdatetime: hash.LocalUpdatetime,
		})
	}
	// 遍历子集
	for _, child := range info.SubCollection {
		childRepo := &Repository{
			ID:          child.ID,
			Title:       child.Title,
			Description: child.Description,
			Version:     "",
			Blueprint:   make([]*Blueprint, 0),
			Repository:  make([]*Repository, 0),
		}
		if err := bfs(ctx, filepath.Join(path, filepath.Clean(child.Title)), childRepo, child.ID, hashMap); err != nil {
			return err
		}
		repo.Repository = append(repo.Repository, childRepo)
	}
	return nil
}
