package migrations

import (
	"context"

	"github.com/codfrm/cago/database/migrate/mongomigrate"
	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/task/producer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func T20240130() *mongomigrate.Migration {
	return &mongomigrate.Migration{
		ID: "20240130",
		Migrate: func(ctx context.Context, db *mongo.Client) error {
			// 蓝图集生成zip到oss
			collections, total, err := collection_repo.Colletcion().FindPage(ctx, httputils.PageRequest{
				Page: 0,
				Size: 100,
			})
			if err != nil {
				return err
			}
			logger.Ctx(ctx).Info("蓝图集数量", zap.Int64("total", total))
			for _, collection := range collections {
				if err := producer.PublishCollectionUpdate(ctx, collection.ID, map[primitive.ObjectID]struct{}{
					collection.ID: {},
				}); err != nil {
					return err
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
