package collection_ctr

import (
	"context"

	api "github.com/dsp2b/dsp2b-go/internal/api/collection"
	"github.com/dsp2b/dsp2b-go/internal/service/collection_svc"
)

type Collection struct {
}

func NewCollection() *Collection {
	return &Collection{}
}

// SubCollection 查询子蓝图集
func (c *Collection) SubCollection(ctx context.Context, req *api.SubCollectionRequest) (*api.SubCollectionResponse, error) {
	return collection_svc.Collection().SubCollection(ctx, req)
}

// GetCollectionBlueprint 查询蓝图
func (c *Collection) GetCollectionBlueprint(ctx context.Context, req *api.GetCollectionBlueprintRequest) (*api.GetCollectionBlueprintResponse, error) {
	return collection_svc.Collection().GetCollectionBlueprint(ctx, req)
}

// Detail 获取蓝图集详情
func (c *Collection) Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error) {
	return collection_svc.Collection().Detail(ctx, req)
}

// UpdateNotify 更新通知
func (c *Collection) UpdateNotify(ctx context.Context, req *api.UpdateNotifyRequest) (*api.UpdateNotifyResponse, error) {
	return collection_svc.Collection().UpdateNotify(ctx, req)
}
