package migrations

import (
	"context"

	"github.com/codfrm/cago/database/migrate/mongomigrate"
	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/model/entity/blueprint_entity"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_repo"
	"github.com/dsp2b/dsp2b-go/internal/service/blueprint_svc"
	"github.com/dsp2b/dsp2b-go/pkg/blueprint"
	"go.uber.org/zap"
)

func T20240201() *mongomigrate.Migration {
	return &mongomigrate.Migration{
		ID: "20240201",
		Migrate: func(ctx context.Context, db *mongo.Client) error {
			// 解析所有蓝图的icon
			for page := 1; ; page++ {
				list, _, err := blueprint_repo.Blueprint().FindPage(ctx, httputils.PageRequest{
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
					decode, err := blueprint.Decode(v.Blueprint)
					if err != nil {
						return err
					}
					if !(decode.Icon0 == 0 && decode.Icon1 == 0 && decode.Icon2 == 0 && decode.Icon3 == 0 && decode.Icon4 == 0 && decode.Icon5 == 0) {
						icons := blueprint_entity.Icons{}
						icons.Icon0, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon0)
						icons.Icon1, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon1)
						icons.Icon2, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon2)
						icons.Icon3, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon3)
						icons.Icon4, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon4)
						icons.Icon5, _ = blueprint_svc.Blueprint().IconInfo(ctx, decode.Icon5)
						icons.Layout = decode.Layout
						if err := v.SetIcons(icons); err != nil {
							return err
						}
					}
					v.Blueprint, err = blueprint.Rename(v.Blueprint, decode.ShortDesc)
					if err != nil {
						return err
					}
					if err := blueprint_repo.Blueprint().Update(ctx, v); err != nil {
						return err
					}
					logger.Ctx(ctx).Info("处理蓝图成功", zap.String("id", v.ID.Hex()))
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
