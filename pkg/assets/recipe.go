package assets

import (
	"encoding/binary"
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
	ReadInt64(f)
	r.ProtoSet.Type = ReadString(f)
	r.TableName = ReadString(f)
	r.Signature = ReadString(f)
	l := ReadInt32(f)
	r.DataArray = make([]RecipeProto, l)
	for i := 0; int32(i) < l; i++ {
		_ = ReadStruct(f, &r.DataArray[i])
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
