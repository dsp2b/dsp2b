package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

type scanBlueprint struct {
	repo          *Repository
	blueprintFile map[string]*Blueprint
	newFile       map[string]os.DirEntry
	modifyFile    map[string]string
}

func newScan(repo *Repository) *scanBlueprint {
	return &scanBlueprint{
		repo:          repo,
		blueprintFile: repo.BlueprintMap(),
		newFile:       make(map[string]os.DirEntry),
		modifyFile:    make(map[string]string),
	}
}

func statusCmd(cmd *cobra.Command, args []string) error {
	// 读取dspb.json
	repo, err := ReadRepository()
	if err != nil {
		return err
	}
	// 比较
	diff := newScan(repo)

	err = diff.Diff(".")
	if err != nil {
		return err
	}

	// 按新文件蓝色 修改文件黄色 删除文件红色输出
	if len(diff.newFile) > 0 {
		for k := range diff.newFile {
			printBlueColor("新蓝图: " + k + "\n")
		}
		fmt.Printf("\n")
	}

	if len(diff.modifyFile) > 0 {
		for k := range diff.modifyFile {
			printYellowColor("修改蓝图: " + k + "\n")
		}

		fmt.Printf("\n")
	}

	if len(diff.blueprintFile) > 0 {
		for k := range diff.blueprintFile {
			printRedColor("删除蓝图: " + k + "\n")
		}
	}

	if err := SaveRepository(repo); err != nil {
		return err
	}

	return nil
}

func (d *scanBlueprint) Diff(path string) error {
	err := d.diff(path)
	if err != nil {
		return err
	}
	// 最后剩下的就是删除的
	return nil
}

func (d *scanBlueprint) Scan(path string, callback func(path string, entry os.DirEntry) error) error {
	err := d.scan(path, callback)
	if err != nil {
		return err
	}
	return nil
}

func (d *scanBlueprint) scan(path string, callback func(path string, entity os.DirEntry) error) error {
	// 遍历所有文件
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, v := range dir {
		if v.IsDir() {
			if err := d.scan(filepath.Join(path, v.Name()), callback); err != nil {
				return err
			}
		} else {
			name := strings.TrimSuffix(v.Name(), ".txt")
			if name == v.Name() {
				continue
			}
			if err := callback(path, v); err != nil {
				logger.Default().Error("文件处理失败", zap.String("file", v.Name()), zap.Error(err))
				return err
			}
		}
	}
	return nil
}

func (d *scanBlueprint) diff(path string) error {
	// 遍历所有文件
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	file := make(map[string]os.DirEntry)
	for _, v := range dir {
		if v.IsDir() {
			if err := d.diff(filepath.Join(path, v.Name())); err != nil {
				return err
			}
		} else {
			name := strings.TrimSuffix(v.Name(), ".txt")
			if name == v.Name() {
				continue
			}
			file[filepath.Join(path, name)] = v
		}
	}
	// 判断是否存在
	for k, v := range file {
		if _, ok := d.blueprintFile[k]; !ok {
			// 不存在则是新文件
			d.newFile[k] = v
		} else {
			// 存在判断修改时间
			info, err := v.Info()
			if err != nil {
				return err
			}
			if info.ModTime().Unix() > d.blueprintFile[k].Updatetime {
				// 本地修改时间大于服务器修改时间, 读取内容判断hash
				hash, err := HashFile(filepath.Join(path, v.Name()))
				if err != nil {
					return err
				}
				if hash != d.blueprintFile[k].Hash {
					d.modifyFile[k] = v.Name()
				} else {
					d.blueprintFile[k].LocalUpdatetime = info.ModTime().Unix()
				}
			}
			delete(d.blueprintFile, k)
		}
	}
	return nil
}
