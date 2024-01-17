package blueprint_svc

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/dsp2b/dsp2b-go/internal/api/blueprint"
	"github.com/dsp2b/dsp2b-go/pkg/assets"
	"github.com/dsp2b/dsp2b-go/pkg/blueprint"
	"go.uber.org/zap"
)

type BlueprintSvc interface {
	// Parse 蓝图解析
	Parse(ctx context.Context, req *api.ParseRequest) (*api.ParseResponse, error)
	// GetRecipePanel 获取配方面板
	GetRecipePanel(ctx context.Context, req *api.GetRecipePanelRequest) (*api.GetRecipePanelResponse, error)
	// Detail 获取蓝图详情
	Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error)
	// List 蓝图列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
}

type blueprintSvc struct {
	ItemProtoMap   map[int32]assets.Proto[assets.ItemProto]
	RecipeProtoMap map[int32]assets.Proto[assets.RecipeProto]
	RecipePanel    api.RecipePanel
}

var defaultBlueprint BlueprintSvc

func Blueprint() BlueprintSvc {
	return defaultBlueprint
}

func InitBlueprint() error {
	b, err := os.ReadFile("data/itemProtoSet.json")
	if err != nil {
		return err
	}
	itemProtoSet := assets.ItemProtoSet{}
	if err := json.Unmarshal(b, &itemProtoSet); err != nil {
		return err
	}
	b, err = os.ReadFile("data/recipeProtoSet.json")
	if err != nil {
		return err
	}
	recipeProtoSet := assets.RecipeProtoSet{}
	if err := json.Unmarshal(b, &recipeProtoSet); err != nil {
		return err
	}
	svc := &blueprintSvc{
		ItemProtoMap:   itemProtoSet.Map(),
		RecipeProtoMap: recipeProtoSet.Map(),
	}
	// 	生成配方面板文件
	panel := api.RecipePanel{}
	for _, v := range recipeProtoSet.DataArray {
		if v.Proto.GridIndex == 0 {
			continue
		}
		num := v.Proto.GridIndex % 1000
		x, y := num/100, num%100-1
		item := api.RecipePanelItem{
			ID:       v.ID,
			ItemID:   v.Proto.Results[0],
			Name:     v.Name,
			IconPath: v.Proto.IconPath,
		}
		if item.IconPath == "" {
			// 找第一个产物的icon
			i, ok := svc.ItemProtoMap[v.Proto.Results[0]]
			if !ok {
				return errors.New("not found icon")
			}
			item.IconPath = i.Proto.IconPath
		}
		item.IconPath = path.Base(item.IconPath)
		if v.Proto.GridIndex > 2000 {
			panel.BuildingPanel[x][y] = item
		} else {
			panel.ThingPanel[x][y] = item
		}
	}
	svc.RecipePanel = panel
	defaultBlueprint = svc
	return nil
}

// Parse 蓝图解析
func (b *blueprintSvc) Parse(ctx context.Context, req *api.ParseRequest) (*api.ParseResponse, error) {
	decode, err := blueprint.Decode(req.Blueprint)
	if err != nil {
		return nil, httputils.NewBadRequestError(-1, err.Error())
	}
	buldingsMap := make(map[int16]*api.Building)
	for _, v := range decode.Buildings {
		if _, ok := buldingsMap[v.ItemId]; ok {
			buldingsMap[v.ItemId].Count += 1
		} else {
			item, ok := b.ItemProtoMap[int32(v.ItemId)]
			if !ok {
				buldingsMap[v.ItemId] = &api.Building{
					ItemId:   v.ItemId,
					Name:     "unknown",
					IconPath: "unknown",
					Count:    1,
				}
				continue
			}
			buldingsMap[v.ItemId] = &api.Building{
				ItemId:   v.ItemId,
				Name:     item.Name,
				IconPath: path.Base(item.Proto.IconPath),
				Count:    1,
			}
		}
	}
	buildings := make([]*api.Building, 0, len(buldingsMap))
	for _, v := range buldingsMap {
		buildings = append(buildings, v)
	}
	// 计算产物
	products, err := b.calcProduct(ctx, decode.Buildings)
	if err != nil {
		return nil, err
	}
	return &api.ParseResponse{
		Blueprint: decode,
		Buildings: buildings,
		Products:  products,
	}, nil
}

// 计算产量
func (b *blueprintSvc) calcProduct(ctx context.Context, buildings []blueprint.Building) ([]*api.Product, error) {
	productMap := make(map[int16]*api.Product)
	for _, v := range buildings {
		// 查看是否有配方
		if v.RecipeId == 0 {
			continue
		}
		recipe, ok := b.RecipeProtoMap[int32(v.RecipeId)]
		if !ok {
			logger.Ctx(ctx).Error("not found recipe", zap.Int16("recipe_id", v.RecipeId))
			continue
		}
		m, p, err := b.recipeSpeed(ctx, recipe)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			if _, ok := productMap[k]; ok {
				productMap[k].Count -= v.Count
			} else {
				v.Count = -v.Count
				productMap[k] = v
			}
		}
		for k, v := range p {
			if _, ok := productMap[k]; ok {
				productMap[k].Count += v.Count
			} else {
				productMap[k] = v
			}
		}
	}
	products := make([]*api.Product, 0, len(productMap))
	for _, v := range productMap {
		v.Count = math.Round(v.Count)
		products = append(products, v)
	}
	return products, nil
}

func (b *blueprintSvc) recipeSpeed(ctx context.Context, recipe assets.Proto[assets.RecipeProto]) (
	map[int16]*api.Product,
	map[int16]*api.Product, error) {
	result := struct {
		ProductMap  map[int16]*api.Product `json:"ProductMap"`
		MaterialMap map[int16]*api.Product `json:"MaterialMap"`
	}{}
	if err := cache.Ctx(ctx).GetOrSet("recipe_speed:"+strconv.Itoa(int(recipe.ID)), func() (interface{}, error) {
		productMap := make(map[int16]*api.Product)
		materialMap := make(map[int16]*api.Product)
		for n, v := range recipe.Proto.Results {
			// 计算速度
			count := recipe.Proto.ResultCounts[n]
			speed := 60 / (float64(recipe.Proto.TimeSpend) / 60) * float64(count)
			if _, ok := productMap[int16(v)]; ok {
				productMap[int16(v)].Count += speed
			} else {
				item, ok := b.ItemProtoMap[v]
				if !ok {
					logger.Ctx(ctx).Error("not found item", zap.Int32("item_id", v))
					continue
				}
				productMap[int16(v)] = &api.Product{
					ItemId:   int16(v),
					Name:     item.Name,
					IconPath: path.Base(item.Proto.IconPath),
					Count:    speed,
				}
			}
		}
		for n, v := range recipe.Proto.Items {
			// 计算速度
			count := recipe.Proto.ItemCounts[n]
			speed := 60 / (float64(recipe.Proto.TimeSpend) / 60) * float64(count)
			if _, ok := materialMap[int16(v)]; ok {
				materialMap[int16(v)].Count += 1
			} else {
				item, ok := b.ItemProtoMap[v]
				if !ok {
					logger.Ctx(ctx).Error("not found item", zap.Int32("item_id", v))
					continue
				}
				materialMap[int16(v)] = &api.Product{
					ItemId:   int16(v),
					Name:     item.Name,
					IconPath: path.Base(item.Proto.IconPath),
					Count:    speed,
				}
			}
		}
		result.MaterialMap = materialMap
		result.ProductMap = productMap
		return result, nil
	}, cache.Expiration(time.Hour)).Scan(&result); err != nil {
		return nil, nil, err
	}
	return result.MaterialMap, result.ProductMap, nil
}

// GetRecipePanel 获取配方面板
func (b *blueprintSvc) GetRecipePanel(ctx context.Context, req *api.GetRecipePanelRequest) (*api.GetRecipePanelResponse, error) {
	return &api.GetRecipePanelResponse{
		RecipePanel: b.RecipePanel,
	}, nil
}

// Detail 获取蓝图详情
func (b *blueprintSvc) Detail(ctx context.Context, req *api.DetailRequest) (*api.DetailResponse, error) {
	return nil, nil
}

// List 蓝图列表
func (b *blueprintSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return nil, nil
}
