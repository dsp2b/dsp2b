package api

import (
	"context"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/dsp2b/dsp2b-go/docs"
	"github.com/dsp2b/dsp2b-go/internal/controller/blueprint_ctr"
	"github.com/dsp2b/dsp2b-go/internal/service/blueprint_svc"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v1")

	rg := r.Group("/")
	if err := blueprint_svc.InitBlueprint(); err != nil {
		return err
	}
	{
		ctr := blueprint_ctr.NewBlueprint()
		rg.Bind(
			ctr.Parse,
			ctr.GetRecipePanel,
		)
	}

	return nil
}
