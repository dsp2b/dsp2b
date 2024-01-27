package resource

import "github.com/codfrm/cago/server/mux"

// ImageThumbnailRequest 缩略图
type ImageThumbnailRequest struct {
	mux.Meta `path:"/image/thumbnail/:width/:height/images/*path" method:"GET"`
	Path     string `uri:"path"`
	Width    uint   `uri:"width"`
	Height   uint   `uri:"height"`
}

type ImageThumbnailResponse struct {
	Redirect string
	Bytes    []byte
}
