package assets

import (
	"encoding/binary"
	"github.com/dsp2b/dsp2b-go/pkg/utils"
	"os"
)

type RecipeProtoSet struct {
	ProtoSet[RecipeProto]
}

func (r *RecipeProtoSet) Load(datFile string) error {
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
	r.ProtoSet.Type = utils.ReadString(f)
	r.TableName = utils.ReadString(f)
	r.Signature = utils.ReadString(f)
	l := utils.ReadInt32(f)
	r.DataArray = make([]RecipeProto, l)
	for i := 0; int32(i) < l; i++ {
		_ = utils.ReadStruct(f, &r.DataArray[i])
	}
	return nil
}

type RecipeProto struct {
	Proto
	Type          ERecipeType
	Handcraft     bool
	Explicit      bool
	TimeSpend     int32   // 速度, 每分钟多少个
	Items         []int32 // 入
	ItemCounts    []int32 // 入数量
	Results       []int32 // 出
	ResultCounts  []int32 // 出数量
	GridIndex     int32
	IconPath      string
	Description   string
	NonProductive bool
}

type ERecipeType int32

const (
	None ERecipeType = iota
	Smelt
	Chemical
	Refine
	Assemble
	Particle
	Exchange
	PhotonStore
	Fractionate
	Research = 15
)
