package assets

import (
	"encoding/binary"
	"github.com/dsp2b/dsp2b-go/pkg/utils"
	"os"
)

type TechProtoSet struct {
	ProtoSet[TechProto]
}

func (i *TechProtoSet) Load(datFile string) error {
	f, err := os.Open(datFile)
	if err != nil {
		return err
	}
	var flag [20]byte
	_ = binary.Read(f, binary.LittleEndian, &flag)
	if flag != [20]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 0, 0,
	} {
		return err
	}
	utils.ReadInt64(f)
	i.ProtoSet.Type = utils.ReadString(f)
	i.TableName = utils.ReadString(f)
	i.Signature = utils.ReadString(f)
	i01 := utils.ReadInt32(f)
	i.DataArray = make([]Proto[TechProto], i01)
	for n := 0; int32(n) < i01; n++ {
		_ = utils.ReadStruct(f, &i.DataArray[n])
	}
	return nil
}

type TechProto struct {
	Desc                  string
	Conclusion            string
	Published             bool
	IsHiddenTech          bool
	IsObsolete            bool
	PreItem               []int32
	Level                 int32
	MaxLevel              int32
	LevelCoef1            int32
	LevelCoef2            int32
	IconPath              string
	IsLabTech             bool
	PreTechs              []int32
	PreTechsImplicit      []int32
	PreTechsMax           bool
	Items                 []int32
	ItemPoints            []int32
	PropertyOverrideItems []int32
	PropertyItemCounts    []int32
	HashNeeded            int64
	UnlockRecipes         []int32
	UnlockFunctions       []int32
	UnlockValues          []float64
	AddItems              []int32
	AddItemCounts         []int32
	Position              Vector2
}
