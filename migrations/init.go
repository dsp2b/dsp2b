package migrations

import (
	"context"
	"time"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/migrate/mongomigrate"
	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

// Migrations 数据库迁移操作
func Migrations(ctx context.Context, cfg *configs.Config) error {
	// 添加分布式锁
	if ok, err := redis.Default().
		SetNX(context.Background(), "migrations", "lock", time.Minute*5).Result(); err != nil {
		logger.Ctx(context.Background()).Error("数据库迁移失败", zap.Error(err))
		return err
	} else if !ok {
		logger.Ctx(context.Background()).Info("数据库迁移已经在执行")
	}
	logger.Ctx(context.Background()).Info("开始执行数据库迁移")
	defer redis.Default().Del(context.Background(), "migrations")
	return run(ctx,
		T20240130,
	)
}

func run(ctx context.Context, fs ...func() *mongomigrate.Migration) error {
	ms := make([]*mongomigrate.Migration, 0, len(fs))
	for _, f := range fs {
		ms = append(ms, f())
	}
	m := mongomigrate.New(ctx, mongo.Default(), ms)
	return m.Migrate(mongomigrate.HasPre())
}
