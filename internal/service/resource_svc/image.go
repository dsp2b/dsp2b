package resource_svc

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/codfrm/cago/pkg/logger"
	api "github.com/dsp2b/dsp2b-go/internal/api/resource"
	"github.com/dsp2b/dsp2b-go/internal/storage"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

type ImageSvc interface {
	// ImageThumbnail 缩略图
	ImageThumbnail(ctx context.Context, req *api.ImageThumbnailRequest) (*api.ImageThumbnailResponse, error)
}

type imageSvc struct {
}

var defaultImage = &imageSvc{}

func Image() ImageSvc {
	return defaultImage
}

// ImageThumbnail 缩略图
func (i *imageSvc) ImageThumbnail(ctx context.Context, req *api.ImageThumbnailRequest) (*api.ImageThumbnailResponse, error) {
	// 读取存储中的图片
	if strings.HasPrefix(req.Path, "images/blueprint") {
		return nil, errors.New("path error")
	}
	thumbnailName := "images" + req.Path
	thumbnailName = path.Dir(thumbnailName) + "/thumbnail/" + path.Base(thumbnailName) + ".jpg"
	// 判断是否有缩略图
	f, err := storage.Default().GetObject(ctx, thumbnailName)
	if err != nil {
		// 缩略图不存在时继续走生成流程
		if !errors.Is(err, storage.ErrNotFound) {
			var e storage.RespError
			if !errors.As(err, &e) || e.StatusCode() != http.StatusNotFound {
				return nil, err
			}
		}
	} else {
		defer f.Close()
		if info, statErr := f.Stat(); statErr == nil && info.Size > 0 {
			data, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}
			return &api.ImageThumbnailResponse{
				Bytes: data,
			}, nil
		}
	}
	r, err := storage.Default().GetObject(ctx, "images"+req.Path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	img, s, err := image.Decode(r)
	// 生成缩略图
	if err != nil {
		return nil, err
	}
	logger.Ctx(ctx).Info("image format", zap.String("path", req.Path), zap.String("format", s))

	// 计算裁剪区域，以保持宽高比
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// 计算裁剪比例和裁剪大小
	var cropWidth, cropHeight int
	targetRatio := float64(req.Width) / float64(req.Height)
	imageRatio := float64(originalWidth) / float64(originalHeight)

	if imageRatio > targetRatio {
		// 图片过宽，需要裁剪宽度
		cropHeight = originalHeight
		cropWidth = int(float64(cropHeight) * targetRatio)
	} else {
		// 图片过高，需要裁剪高度
		cropWidth = originalWidth
		cropHeight = int(float64(cropWidth) / targetRatio)
	}

	// 计算裁剪的起始点（以居中裁剪）
	startX := (originalWidth - cropWidth) / 2
	startY := (originalHeight - cropHeight) / 2

	// 裁剪图片
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(startX, startY, startX+cropWidth, startY+cropHeight))

	// 缩放图片
	resizedImg := resize.Resize(req.Width, req.Height, croppedImg, resize.NearestNeighbor)

	buf := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buf, resizedImg, &jpeg.Options{
		Quality: 80,
	}); err != nil {
		return nil, err
	}
	// 写入存储，作为缩略图缓存
	if err := storage.Default().PutObject(ctx, thumbnailName, bytes.NewReader(buf.Bytes()), int64(buf.Len())); err != nil {
		return nil, err
	}
	return &api.ImageThumbnailResponse{
		Bytes: buf.Bytes(),
	}, nil
}
