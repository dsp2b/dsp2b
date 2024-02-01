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

	"github.com/dsp2b/dsp2b-go/internal/model/entity/blueprint_entity"

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
	// Create 创建蓝图
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
	// IconInfo 获取icon信息
	IconInfo(ctx context.Context, itemId int) (*blueprint_entity.IconInfo, error)
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

func InitBlueprint(itemProtoSetPath, recipeProtoSetPath string) error {
	b, err := os.ReadFile(itemProtoSetPath)
	if err != nil {
		return err
	}
	itemProtoSet := assets.ItemProtoSet{}
	if err := json.Unmarshal(b, &itemProtoSet); err != nil {
		return err
	}
	b, err = os.ReadFile(recipeProtoSetPath)
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
	for k, v := range svc.ItemProtoMap {
		v.Proto.IconPath = path.Base(v.Proto.IconPath)
		svc.ItemProtoMap[k] = v
	}
	// 	生成配方面板文件
	panel := api.RecipePanel{}
	for k, v := range svc.RecipeProtoMap {
		v.Proto.IconPath = path.Base(v.Proto.IconPath)
		svc.RecipeProtoMap[k] = v
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
				IconPath: item.Proto.IconPath,
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

// SpeedMultiplier 速度倍率
type SpeedMultiplier struct {
	Name string
	// 速度倍率
	Multiplier float64
	// 增产速度
	Productivity float64
}

var speedMultiplierMap = map[int16]SpeedMultiplier{
	2302: {"电弧熔炉", 1, 0},
	2315: {"位面熔炉", 2, 0},
	2319: {"熔炉 Mk.III", 3, 0},

	2303: {"制作台Mk.Ⅰ", 0.75, 0},
	2304: {"制作台Mk.Ⅱ", 1, 0},
	2305: {"制作台Mk.Ⅲ", 1.5, 0},
	2318: {"重组式制造台", 3, 0},

	2901: {"矩阵研究站", 1, 0},
	2902: {"自演化研究站", 3, 0},

	2309: {"化工厂", 1, 0},
	2317: {"量子化工厂", 2, 0},

	1141: {"增产剂 Mk.I", 1.25, 1.125},
	1142: {"增产剂 Mk.II", 1.5, 1.2},
	1143: {"增产剂 Mk.III", 2, 1.25},
}

// 计算产量
func (b *blueprintSvc) calcProduct(ctx context.Context, buildings []blueprint.Building) ([]*api.Product, error) {
	productMap := make(map[int16]*api.Product)
	// 判断是否有喷涂机, 有涂装直接算上三级增产剂
	increase := false
	for _, v := range buildings {
		if v.ItemId == 2313 {
			increase = true
			break
		}
	}
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
		var (
			// 生产速度
			pSpeedMultiplier float64 = 1
			// 消费速度
			mSpeedMultiplier float64 = 1
		)
		speedMultiplierValue, ok := speedMultiplierMap[v.ItemId]
		if ok {
			pSpeedMultiplier = speedMultiplierValue.Multiplier
			mSpeedMultiplier = speedMultiplierValue.Multiplier
		}
		if increase {
			increaseMultiplier := speedMultiplierMap[1143]
			// 判断增产还是增速
			if len(v.Parameters) > 0 {
				if v.Parameters[0] == 1 {
					// 加速
					pSpeedMultiplier = pSpeedMultiplier * increaseMultiplier.Multiplier
					mSpeedMultiplier = mSpeedMultiplier * increaseMultiplier.Multiplier
				} else {
					// 增产
					pSpeedMultiplier = pSpeedMultiplier * increaseMultiplier.Productivity
				}
			}
		}
		for k, v := range m {
			if _, ok := productMap[k]; ok {
				productMap[k].Count -= v.Count * mSpeedMultiplier
			} else {
				v.Count = -v.Count * mSpeedMultiplier
				productMap[k] = v
			}
		}
		for k, v := range p {
			if _, ok := productMap[k]; ok {
				productMap[k].Count += v.Count * pSpeedMultiplier
			} else {
				v.Count = v.Count * pSpeedMultiplier
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
					IconPath: item.Proto.IconPath,
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
					IconPath: item.Proto.IconPath,
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

// Create 创建蓝图
func (b *blueprintSvc) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return nil, nil
}

func (b *blueprintSvc) IconInfo(ctx context.Context, itemId int) (*blueprint_entity.IconInfo, error) {
	item, ok := b.ItemProtoMap[int32(itemId)]
	if !ok {
		return nil, errors.New("not found item")
	}
	return &blueprint_entity.IconInfo{
		ItemID:   item.ID,
		Name:     item.Name,
		IconPath: item.Proto.IconPath,
	}, nil
}
