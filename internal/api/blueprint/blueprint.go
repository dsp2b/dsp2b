package blueprint

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/dsp2b/dsp2b-go/pkg/blueprint"
)

// ParseRequest 蓝图解析
type ParseRequest struct {
	mux.Meta  `path:"/blueprint/parse" method:"POST""`
	Blueprint string `form:"blueprint" binding:"required"`
}

type ParseResponse struct {
	Blueprint blueprint.Blueprint `json:"blueprint"`
}
