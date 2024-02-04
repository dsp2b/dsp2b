package blueprint_ctr

import (
	"context"

	api "github.com/dsp2b/dsp2b-go/internal/api/blueprint"
	"github.com/dsp2b/dsp2b-go/internal/service/blueprint_svc"
)

type Blueprint struct {
}

func NewBlueprint() *Blueprint {
	return &Blueprint{}
}

// Parse 蓝图解析
func (b *Blueprint) Parse(ctx context.Context, req *api.ParseRequest) (*api.ParseResponse, error) {
	return blueprint_svc.Blueprint().Parse(ctx, req)
}

// GetRecipePanel 获取配方面板
func (b *Blueprint) GetRecipePanel(ctx context.Context, req *api.GetRecipePanelRequest) (*api.GetRecipePanelResponse, error) {
	return blueprint_svc.Blueprint().GetRecipePanel(ctx, req)
}

// Detail 获取蓝图详情
func (b *Blueprint) Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error) {
	return blueprint_svc.Blueprint().Detail(ctx, req)
}

// List 蓝图列表
func (b *Blueprint) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return blueprint_svc.Blueprint().List(ctx, req)
}

// Create 创建蓝图
func (b *Blueprint) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return blueprint_svc.Blueprint().Create(ctx, req)
}

// ReplaceBlueprint 替换蓝图配方和建筑
func (b *Blueprint) ReplaceBlueprint(ctx context.Context, req *api.ReplaceBlueprintRequest) (*api.ReplaceBlueprintResponse, error) {
	return blueprint_svc.Blueprint().ReplaceBlueprint(ctx, req)
}
