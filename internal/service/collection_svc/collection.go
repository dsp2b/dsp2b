package collection_svc

import (
	"context"
	"errors"

	"github.com/dsp2b/dsp2b-go/internal/task/producer"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/blueprint_repo"
	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"

	api "github.com/dsp2b/dsp2b-go/internal/api/collection"
)

type CollectionSvc interface {
	// SubCollection 查询子蓝图集
	SubCollection(ctx context.Context, req *api.SubCollectionRequest) (*api.SubCollectionResponse, error)
	// GetCollectionBlueprint 查询蓝图
	GetCollectionBlueprint(ctx context.Context, req *api.GetCollectionBlueprintRequest) (*api.GetCollectionBlueprintResponse, error)
	// Detail 获取蓝图集详情
	Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error)
	// UpdateNotify 更新通知
	UpdateNotify(ctx context.Context, req *api.UpdateNotifyRequest) (*api.UpdateNotifyResponse, error)
}

type collectionSvc struct {
}

var defaultCollection = &collectionSvc{}

func Collection() CollectionSvc {
	return defaultCollection
}

// SubCollection 查询子蓝图集
func (c *collectionSvc) SubCollection(ctx context.Context, req *api.SubCollectionRequest) (*api.SubCollectionResponse, error) {
	list, err := collection_repo.Colletcion().FindByParent(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	resp := &api.SubCollectionResponse{
		Collection: make([]api.Collection, 0),
	}
	for _, v := range list {
		resp.Collection = append(resp.Collection, api.Collection{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
		})
	}
	return resp, nil
}

// GetCollectionBlueprint 查询蓝图
func (c *collectionSvc) GetCollectionBlueprint(ctx context.Context, req *api.GetCollectionBlueprintRequest) (*api.GetCollectionBlueprintResponse, error) {
	list, err := blueprint_collection_repo.BlueprintCollection().FindByCollection(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	resp := &api.GetCollectionBlueprintResponse{
		Blueprint: make([]api.GetCollectionBlueprintItem, 0),
	}
	for _, v := range list {
		blueprint, err := blueprint_repo.Blueprint().Find(ctx, v.BlueprintID)
		if err != nil {
			return nil, err
		}
		resp.Blueprint = append(resp.Blueprint, api.GetCollectionBlueprintItem{
			Title:       blueprint.Title,
			Description: blueprint.Description,
			Blueprint:   blueprint.Blueprint,
			Createtime:  blueprint.Createtime.Unix(),
			Updatetime:  blueprint.Updatetime.Unix(),
		})
	}
	return resp, nil
}

var (
	ErrNotFound = errors.New("collection not found")
)

// Detail 获取蓝图集详情
func (c *collectionSvc) Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error) {
	collection, err := collection_repo.Colletcion().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, ErrNotFound
	}
	return &api.DetailResponse{
		ID:          collection.ID,
		Title:       collection.Title,
		Description: collection.Description,
		ParentID:    collection.ParentID,
	}, nil
}

// UpdateNotify 更新通知
func (c *collectionSvc) UpdateNotify(ctx context.Context, req *api.UpdateNotifyRequest) (*api.UpdateNotifyResponse, error) {
	if err := producer.PublishCollectionUpdate(ctx, req.ID, map[primitive.ObjectID]struct{}{
		req.ID: {},
	}); err != nil {
		return nil, err
	}
	return &api.UpdateNotifyResponse{}, nil
}
