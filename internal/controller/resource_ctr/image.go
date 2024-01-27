package resource_ctr

import (
	"net/http"
	"strconv"

	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/dsp2b/dsp2b-go/internal/api/resource"
	"github.com/dsp2b/dsp2b-go/internal/service/resource_svc"
	"github.com/gin-gonic/gin"
)

type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

// ImageThumbnail 缩略图
func (i *Image) ImageThumbnail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		w, _ := strconv.Atoi(ctx.Param("width"))
		if w < 100 || w > 500 {
			w = 280
		}
		h, _ := strconv.Atoi(ctx.Param("height"))
		if h < 100 || h > 500 {
			h = 200
		}
		req := &api.ImageThumbnailRequest{
			Path:   ctx.Param("path"),
			Width:  uint(w),
			Height: uint(h),
		}
		resp, err := resource_svc.Image().ImageThumbnail(ctx, req)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if resp.Redirect != "" {
			ctx.Redirect(http.StatusMovedPermanently, resp.Redirect)
			return
		}
		// 缓存1天
		ctx.Header("Cache-Control", "max-age=86400")
		ctx.Data(http.StatusOK, "image/jpeg", resp.Bytes)
	}
}
