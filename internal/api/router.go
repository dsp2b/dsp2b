package api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	_ "github.com/dsp2b/dsp2b-go/docs"
	"github.com/dsp2b/dsp2b-go/internal/controller/blueprint_ctr"
	"github.com/dsp2b/dsp2b-go/internal/controller/collection_ctr"
	"github.com/dsp2b/dsp2b-go/internal/controller/resource_ctr"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v1")
	rpc := root.Group("/rpc")
	rg := r.Group("/")
	rpcg := rpc.Group("/")
	{
		ctr := blueprint_ctr.NewBlueprint()
		rpcg.Bind(
			ctr.Parse,
			ctr.GetRecipePanel,
		)
	}
	{
		ctr := collection_ctr.NewCollection()
		rpcg.Bind(
			ctr.UpdateNotify,
		)
	}

	{
		ctr := resource_ctr.NewImage()
		rg.GET("/image/thumbnail/:width/:height/images/*path", ctr.ImageThumbnail())
	}
	return nil
}
