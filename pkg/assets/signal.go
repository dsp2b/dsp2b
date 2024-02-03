package assets

import (
	"encoding/binary"
	"os"

	"github.com/dsp2b/dsp2b-go/pkg/utils"
)

type SignalProtoSet struct {
	ProtoSet[SignalProto]
}

func (i *SignalProtoSet) Load(datFile string) error {
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
	i.DataArray = make([]Proto[SignalProto], i01)
	for n := 0; int32(n) < i01; n++ {
		_ = utils.ReadStruct(f, &i.DataArray[n])
	}
	return nil
}

type SignalProto struct {
	GridIndex int32
	IconPath  string
}
