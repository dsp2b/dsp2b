package collection

import (
	"github.com/codfrm/cago/server/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadRequest 下载蓝图zip包
type DownloadRequest struct {
	mux.Meta `path:"/collection/:id/download" method:"GET"`
	ID       primitive.ObjectID `path:"id"`
}

// DownloadResponse 下载蓝图zip包
type DownloadResponse struct {
}

type Collection struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
}

// SubCollectionRequest 查询子蓝图集
type SubCollectionRequest struct {
	mux.Meta `path:"/collection/:id/sub" method:"GET"`
	ID       primitive.ObjectID `path:"id"`
}

type SubCollectionResponse struct {
	Collection []Collection `json:"collection"`
}

// GetCollectionBlueprintRequest 查询蓝图
type GetCollectionBlueprintRequest struct {
	mux.Meta `path:"/collection/:id/blueprint" method:"GET"`
	ID       primitive.ObjectID `path:"id"`
}

type GetCollectionBlueprintItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Blueprint   string `json:"blueprint"`
	Createtime  int64  `json:"createtime"`
	Updatetime  int64  `json:"updatetime"`
}

type GetCollectionBlueprintResponse struct {
	Blueprint []GetCollectionBlueprintItem `json:"blueprint"`
}

type DetailRequest struct {
	mux.Meta `path:"/collection/:id" method:"GET"`
	ID       primitive.ObjectID `path:"id"`
}

type DetailResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
}
