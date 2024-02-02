package command

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	api "github.com/dsp2b/dsp2b-go/internal/api/blueprint"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func push(cmd *cobra.Command, args []string) error {
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
	apiClient, err := NewApiClient()
	if err != nil {
		return err
	}
	repoMap := repo.RepositoryMap()

	ctx := context.Background()
	// 暂时只处理新蓝图
	for k, v := range diff.newFile {
		dir := filepath.Dir(k)
		collection, ok := repoMap[dir]
		if !ok {
			if dir == "." {
				collection = repo
			} else {
				logger.Ctx(ctx).Error("蓝图集不存在, 请手动创建蓝图集, 再执行dsp pullCmd", zap.String("dir", dir))
				break
			}
		}
		info, err := v.Info()
		if err != nil {
			logger.Ctx(ctx).Error("读取蓝图文件信息失败", zap.Error(err), zap.String("file", k))
			break
		}
		data, err := os.ReadFile(k + ".txt")
		if err != nil {
			logger.Ctx(ctx).Error("读取蓝图文件失败", zap.Error(err), zap.String("file", k))
		}
		hash, err := Hash(data)
		if err != nil {
			logger.Ctx(ctx).Error("计算蓝图hash失败", zap.Error(err), zap.String("file", k))
			break
		}
		parse, err := apiClient.ParseBlueprint(ctx, &api.ParseRequest{
			Blueprint: string(data),
		})
		if err != nil {
			logger.Ctx(ctx).Error("解析蓝图文件失败", zap.Error(err), zap.String("file", k))
			break
		}
		tags := make([]*api.Tag, 0)
		for _, v := range parse.Products {
			if v.Count > 0 {
				tags = append(tags, &api.Tag{
					ItemId:   int(v.ItemId),
					Name:     v.Name,
					IconPath: v.IconPath,
				})
			}
		}
		resp, err := apiClient.PostBlueprint(ctx, &api.CreateRequest{
			Blueprint:   string(data),
			Collections: []primitive.ObjectID{collection.ID},
			Title:       parse.Blueprint.ShortDesc,
			Description: parse.Blueprint.Desc,
			Products:    parse.Products,
			PicList:     nil,
			Tags:        tags,
		})
		if err != nil {
			logger.Ctx(ctx).Error("发布蓝图失败", zap.Error(err), zap.String("file", k))
			break
		}
		logger.Ctx(ctx).Info("发布蓝图成功", zap.String("file", k), zap.String("id", resp.ID.Hex()))
		collection.Blueprint = append(collection.Blueprint, &Blueprint{
			ID:              resp.ID,
			Title:           parse.Blueprint.ShortDesc,
			Description:     parse.Blueprint.Desc,
			Hash:            hash,
			Updatetime:      time.Now().Unix(),
			LocalUpdatetime: info.ModTime().Unix(),
		})
	}

	if err := SaveRepository(repo); err != nil {
		return err
	}

	return nil
}
