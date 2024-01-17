package blueprint

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/dsp2b/dsp2b-go/pkg/blueprint"
)

type Blueprint struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Blueprint   string `json:"blueprint"`
}

type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ListRequest 蓝图列表
type ListRequest struct {
	mux.Meta              `path:"/blueprint" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

type DetailRequest struct {
	mux.Meta `path:"/blueprint/:id" method:"GET"`
	ID       string `uri:"id"`
}

type DetailResponse struct {
	Blueprint `json:",inline"`
}

// ParseRequest 蓝图解析
type ParseRequest struct {
	mux.Meta  `path:"/blueprint/parse" method:"POST"`
	Blueprint string `form:"blueprint" binding:"required"`
}

type Building struct {
	ItemId   int16  `json:"item_id"`
	Name     string `json:"name"`
	IconPath string `json:"icon_path"`
	Count    int    `json:"count"`
}

type ParseResponse struct {
	Blueprint blueprint.Blueprint `json:"blueprint"`
	Buildings []*Building         `json:"buildings"`
	Products  []*Product          `json:"products"`
}

type Product struct {
	ItemId   int16   `json:"item_id"`
	Name     string  `json:"name"`
	IconPath string  `json:"icon_path"`
	Count    float64 `json:"count"`
}

// GetRecipePanelRequest 获取配方面板
type GetRecipePanelRequest struct {
	mux.Meta `path:"/blueprint/recipe_panel" method:"GET"`
}

type GetRecipePanelResponse struct {
	RecipePanel `json:",inline"`
}

type RecipePanel struct {
	ThingPanel    [8][14]RecipePanelItem `json:"thing_panel"`
	BuildingPanel [8][14]RecipePanelItem `json:"building_panel"`
}

type RecipePanelItem struct {
	ID       int32  `json:"id"`
	ItemID   int32  `json:"item_id"`
	Name     string `json:"name"`
	IconPath string `json:"icon_path"`
}
