package assets

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path"
	"testing"
)

type AssetsXML struct {
	XMLName xml.Name `xml:"Assets"`
	Asset   []*Asset `xml:"Asset"`
}

type Asset struct {
	Name      string `xml:"ShortDesc"`
	Container string `xml:"Container"`
}

func TestRecipeProtoSet_Load(t *testing.T) {
	recipe := &RecipeProtoSet{}
	_ = recipe.Load("RecipeProtoSet.dat")
	b, _ := json.MarshalIndent(recipe, "", "  ")
	_ = os.WriteFile("data/RecipeProtoSet.json", b, 0644)

	item := &ItemProtoSet{}
	_ = item.Load("ItemProtoSet.dat")
	b, _ = json.MarshalIndent(item, "", "  ")
	_ = os.WriteFile("data/ItemProtoSet.json", b, 0644)
	// 提取图片
	_ = os.MkdirAll("data/icons/item_recipe", 0755)
	for _, v := range item.DataArray {
		// 从data/dsp/Texture2D中提取图片
		if v.Proto.IconPath == "" {
			continue
		}
		filename := path.Base(v.Proto.IconPath)
		oldName := "data/dsp/Texture2D/" + filename + ".png"
		// copy文件
		err := os.Link(oldName,
			"data/icons/item_recipe/"+filename+".png")
		if err != nil {
			if !os.IsExist(err) {
				t.Fatalf("link err: %v", err)
			}
		}
	}
	for _, v := range recipe.DataArray {
		// 从data/dsp/Texture2D中提取图片
		if v.Proto.IconPath == "" {
			continue
		}
		filename := path.Base(v.Proto.IconPath)
		oldName := "data/dsp/Texture2D/" + filename + ".png"
		// copy文件
		err := os.Link(oldName,
			"data/icons/item_recipe/"+filename+".png")
		if err != nil {
			if !os.IsExist(err) {
				t.Fatalf("link err: %v", err)
			}
		}
	}
}
