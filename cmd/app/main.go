package main

import (
	"context"
	"github.com/codfrm/cago/pkg/component"
	"github.com/dsp2b/dsp2b-go/internal/api"
	"log"

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
	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(component.Cache()).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
