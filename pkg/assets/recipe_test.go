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
	item := &ItemProtoSet{}
	_ = item.Load("ItemProtoSet.dat")

	tech := &TechProtoSet{}
	_ = tech.Load("TechProtoSet.dat")

	signal := &SignalProtoSet{}
	_ = signal.Load("SignalProtoSet.dat")

	// 提取图片
	_ = os.MkdirAll("data/icons/item_recipe", 0755)
	_ = os.MkdirAll("data/icons/tech", 0755)
	_ = os.MkdirAll("data/icons/signal", 0755)
	mvIcon := func(iconPath string) string {
		filename := filepath.Join("data/dsp/asset/", strings.ToLower(iconPath)+".png")
		// copy文件
		newFilename := strings.ToLower(utils.ToPathUnderline(iconPath))
		err := os.Link(filename, filepath.Join("data", newFilename+".png"))
		if err != nil {
			if !os.IsExist(err) {
				t.Fatalf("link err: %v", err)
			}
		}
		return newFilename
	}
	for k, v := range item.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		item.DataArray[k].Proto.IconPath = mvIcon(item.DataArray[k].Proto.IconPath)
	}
	for k, v := range recipe.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		recipe.DataArray[k].Proto.IconPath = mvIcon(recipe.DataArray[k].Proto.IconPath)
	}
	for k, v := range tech.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		tech.DataArray[k].Proto.IconPath = mvIcon(tech.DataArray[k].Proto.IconPath)
	}
	for k, v := range signal.DataArray {
		if v.Proto.IconPath == "" {
			continue
		}
		signal.DataArray[k].Proto.IconPath = mvIcon(signal.DataArray[k].Proto.IconPath)
	}
	b, _ := json.MarshalIndent(recipe, "", "  ")
	_ = os.WriteFile("data/RecipeProtoSet.json", b, 0644)

	b, _ = json.MarshalIndent(item, "", "  ")
	_ = os.WriteFile("data/ItemProtoSet.json", b, 0644)

	b, _ = json.MarshalIndent(tech, "", "  ")
	_ = os.WriteFile("data/TechProtoSet.json", b, 0644)

	b, _ = json.MarshalIndent(signal, "", "  ")
	_ = os.WriteFile("data/SignalProtoSet.json", b, 0644)
}
