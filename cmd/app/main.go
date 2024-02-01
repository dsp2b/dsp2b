package main

import (
	"context"
	"log"

	"github.com/dsp2b/dsp2b-go/internal/service/blueprint_svc"

	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/task/consumer"
	"github.com/dsp2b/dsp2b-go/migrations"

	"github.com/codfrm/cago/pkg/oss"

	"github.com/codfrm/cago/pkg/component"
	"github.com/dsp2b/dsp2b-go/internal/api"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/mux"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("dsp2b")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	{
		blueprint_repo.RegisterBlueprint(blueprint_repo.NewBlueprint())
		collection_repo.RegisterColletcion(collection_repo.NewColletcion())
		blueprint_collection_repo.RegisterBlueprintCollection(blueprint_collection_repo.NewBlueprintCollection())

		if err := blueprint_svc.InitBlueprint(
			"./data/itemProtoSet.json",
			"./data/recipeProtoSet.json",
		); err != nil {
			log.Fatalf("init blueprint err: %v", err)
		}
	}
	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(component.Cache()).
		Registry(component.Redis()).
		Registry(component.Mongo()).
		Registry(component.Broker()).
		Registry(cago.FuncComponent(consumer.Consumer)).
		Registry(cago.FuncComponent(oss.OSS)).
		Registry(cago.FuncComponent(migrations.Migrations)).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
