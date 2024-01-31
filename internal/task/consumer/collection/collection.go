package collection

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/oss"
	api "github.com/dsp2b/dsp2b-go/internal/api/collection"
	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"
	"github.com/dsp2b/dsp2b-go/internal/service/collection_svc"
	"github.com/dsp2b/dsp2b-go/internal/task/producer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Collection struct {
}

func (c *Collection) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeCollectionUpdate(ctx, c.Update); err != nil {
		return err
	}
	return nil
}

func (c *Collection) Update(ctx context.Context, id primitive.ObjectID, exist map[primitive.ObjectID]struct{}) error {
	// 打包蓝图集zip文件上传到oss
	collection, err := collection_svc.Collection().Detail(ctx, &api.DetailRequest{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, collection_svc.ErrNotFound) {
			logger.Ctx(ctx).Error("collection not found", zap.String("id", id.Hex()))
			return nil
		}
		return err
	}
	buf := bytes.NewBuffer(nil)
	filename := "blueprint/collection/" + collection.ID.Hex() +
		"/" + filepath.Clean(collection.Title) + ".zip"
	if err := func() error {
		w := zip.NewWriter(buf)
		defer w.Close()
		if len(collection.Description) > 3000 {
			collection.Description = collection.Description[:3000]
		}
		if err := w.SetComment(collection.Description); err != nil {
			return err
		}
		err = c.writeZip(ctx, w, id, "", map[primitive.ObjectID]struct{}{
			id: {},
		})
		if err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return err
	}
	if err := oss.DefaultBucket().PutObject(ctx, filename, buf); err != nil {
		return err
	}
	// 更新蓝图集下载地址
	if err := collection_repo.Colletcion().UpdateDownloadFile(ctx, id, filename); err != nil {
		return err
	}
	logger.Ctx(ctx).Info("update collection download file success", zap.String("filename", filename))
	// 更新父蓝图集
	if !collection.ParentID.IsZero() {
		if _, ok := exist[collection.ParentID]; ok {
			// 存在闭环
			logger.Ctx(ctx).Error("collection update exist loop", zap.String("id", id.Hex()))
			return nil
		}
		exist[collection.ParentID] = struct{}{}
		if err := producer.PublishCollectionUpdate(ctx, collection.ParentID, exist); err != nil {
			return err
		}
	}
	return nil
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
	for _, v := range blueprint.Blueprint {
		fh := &zip.FileHeader{
			Name:     filepath.Join(pathname, filepath.Clean(v.Title)+".txt"),
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
		err := c.writeZip(ctx, w, v.ID, filepath.Join(pathname, filepath.Clean(v.Title)), subCollectionMap)
		if err != nil {
			return err
		}
	}
	return nil
}
