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
