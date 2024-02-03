package assets

import (
	"encoding/json"
	"encoding/xml"
	"github.com/dsp2b/dsp2b-go/pkg/utils"
	"os"
	"path/filepath"
	"strings"
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

	tech := &TechProtoSet{}
	_ = tech.Load("TechProtoSet.dat")
	b, _ = json.MarshalIndent(tech, "", "  ")
	_ = os.WriteFile("data/TechProtoSet.json", b, 0644)

	signal := &SignalProtoSet{}
	_ = signal.Load("SignalProtoSet.dat")
	b, _ = json.MarshalIndent(signal, "", "  ")
	_ = os.WriteFile("data/SignalProtoSet.json", b, 0644)

	// 提取图片
	_ = os.MkdirAll("data/icons/item_recipe", 0755)
	mvIcon := func(iconPath string) {
		filename := filepath.Join("data/dsp/asset/", strings.ToLower(iconPath)+".png")
		// copy文件
		err := os.Link(filename, filepath.Join("data", strings.ToLower(utils.ToPathUnderline(iconPath)+".png")))
		if err != nil {
			if !os.IsExist(err) {
				t.Fatalf("link err: %v", err)
			}
		}
	}
	for _, v := range item.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		mvIcon(v.Proto.IconPath)
	}
	for _, v := range recipe.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		mvIcon(v.Proto.IconPath)
	}
	for _, v := range tech.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		mvIcon(v.Proto.IconPath)
	}
	for _, v := range signal.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		mvIcon(v.Proto.IconPath)
	}
}
