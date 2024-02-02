package migrations

import (
	"context"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/codfrm/cago/database/migrate/mongomigrate"
	"github.com/codfrm/cago/database/mongo"
)

func T20240130() *mongomigrate.Migration {
	return &mongomigrate.Migration{
		ID: "20240130",
		Migrate: func(ctx context.Context, db *mongo.Client) error {
			// 构建根蓝图集
			rootMap := make(map[primitive.ObjectID]primitive.ObjectID)
			for page := 1; ; page++ {
				list, _, err := blueprint_collection_repo.BlueprintCollection().FindPage(ctx, httputils.PageRequest{
					Page: page,
					Size: 20,
				})
				if err != nil {
					return err
				}
				if len(list) == 0 {
					break
				}
				for _, v := range list {
					root, err := utils.RootCollection(ctx, v.CollectionID)
					if err != nil {
						return err
					}
					rootMap[v.CollectionID] = root
					v.RootCollectionID = root
					if err := blueprint_collection_repo.BlueprintCollection().Update(ctx, v); err != nil {
						return err
					}
					logger.Ctx(ctx).Info("处理蓝图集映射成功", zap.String("id", v.ID.Hex()))
				}
			}
			return nil
		},
		Rollback: func(ctx context.Context, db *mongo.Client) error {
			// 回滚数据库
			return nil
		},
	}
}
