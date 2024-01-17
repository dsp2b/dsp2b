package collection_ctr

import (
	"archive/zip"
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/dsp2b/dsp2b-go/internal/api/collection"
	"github.com/dsp2b/dsp2b-go/internal/service/collection_svc"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
}

func NewCollection() *Collection {
	return &Collection{}
}

// Download 下载蓝图zip包
func (c *Collection) Download() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		collection, err := collection_svc.Collection().Detail(ctx, &api.DetailRequest{
			ID: id,
		})
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		// 文件名header
		ctx.Header("Content-Disposition", "attachment; filename="+filepath.Clean(collection.Title)+".zip")
		w := zip.NewWriter(ctx.Writer)
		defer w.Close()
		err = c.writeZip(ctx, w, id, "", map[primitive.ObjectID]struct{}{
			id: {},
		})
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
	}
}

func (c *Collection) writeZip(ctx context.Context, w *zip.Writer, id primitive.ObjectID,
	pathname string, subCollectionMap map[primitive.ObjectID]struct{}) error {
	// 查询蓝图子集
	subCollection, err := collection_svc.Collection().SubCollection(ctx, &api.SubCollectionRequest{
		ID: id,
	})
	if err != nil {
		return err
	}
	// 查询蓝图
	blueprint, err := collection_svc.Collection().GetCollectionBlueprint(ctx, &api.GetCollectionBlueprintRequest{
		ID: id,
	})
	if err != nil {
		return err
	}
	// 写入蓝图子集
	for _, v := range blueprint.Blueprint {
		fh := &zip.FileHeader{
			Name:     filepath.Join(pathname, strings.ReplaceAll(v.Title, "/", " ")+".txt"),
			Comment:  v.Description,
			Method:   zip.Deflate,
			Modified: time.Unix(v.Updatetime, 0),
		}
		f, err := w.CreateHeader(fh)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(v.Blueprint))
		if err != nil {
			return err
		}
	}
	// 遍历子集
	for _, v := range subCollection.Collection {
		if _, ok := subCollectionMap[v.ID]; ok {
			continue
		}
		subCollectionMap[v.ID] = struct{}{}
		err := c.writeZip(ctx, w, v.ID, filepath.Join(pathname, v.Title), subCollectionMap)
		if err != nil {
			return err
		}
	}
	return nil
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
