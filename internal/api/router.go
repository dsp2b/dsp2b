package api

import (
	"context"
	"github.com/dsp2b/dsp2b-go/internal/controller/collection_ctr"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"

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
	{
		blueprint_repo.RegisterBlueprint(blueprint_repo.NewBlueprint())
		collection_repo.RegisterColletcion(collection_repo.NewColletcion())
		blueprint_collection_repo.RegisterBlueprintCollection(blueprint_collection_repo.NewBlueprintCollection())
	}

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
	{
		ctr := collection_ctr.NewCollection()
		rg.GET("/collection/:id/download", ctr.Download())
	}

	return nil
}
