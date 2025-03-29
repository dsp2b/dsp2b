package blueprint

import (
	"io"

	"github.com/dsp2b/dsp2b-go/pkg/utils"
)

type Area struct {
	Index              int8
	ParentIndex        int8
	TropicAnchor       int16
	AreaSegments       int16
	AnchorLocalOffsetX int16
	AnchorLocalOffsetY int16
	Width              int16
	Height             int16
}

type Building struct {
	Version          int32 // 版本 =-100时，增加了Tilt字段
	Index            int32
	AreaIndex        int8
	LocalOffsetX     float32
	LocalOffsetY     float32
	LocalOffsetZ     float32
	LocalOffsetX2    float32
	LocalOffsetY2    float32
	LocalOffsetZ2    float32
	Yaw              float32
	Yaw2             float32
	Tilt             float32 `binary:"version=-100"`
	Tilt2            float32 `binary:"version=-101"`
	Pitch            float32
	Pitch2           float32
	ItemId           int16
	ModelIndex       int16
	OutputObj        int32
	InputObj         int32
	OutputToSlot     int8
	InputFromSlot    int8
	OutputFromSlot   int8
	InputToSlot      int8
	OutputOffset     int8
	InputOffset      int8
	RecipeId         int16
	FilterId         int16
	Parameters       []int32 `binary:"-"`
	TempOutputObjIdx int32
	TempInputObjIdx  int32
	Content          string
}

// Parse BlueprintBuilding.Export
func (b *Building) Parse(r io.Reader) error {
	version := utils.ReadInt32(r)
	b.Version = version
	if version <= -102 {
		b.Index = utils.ReadInt32(r)
		b.ItemId = utils.ReadInt16(r)
		b.ModelIndex = utils.ReadInt16(r)
		b.AreaIndex = utils.ReadInt8(r)
		b.LocalOffsetX = float32(utils.ReadInt32(r))
		b.LocalOffsetY = float32(utils.ReadInt32(r))
		b.LocalOffsetZ = float32(utils.ReadInt32(r))
		b.Yaw = float32(utils.ReadInt32(r))
		if b.ItemId > 2000 && b.ItemId < 2010 {
			b.Tilt = float32(utils.ReadInt32(r))
			b.Pitch = 0
			b.LocalOffsetX2 = b.LocalOffsetX
			b.LocalOffsetY2 = b.LocalOffsetY
			b.LocalOffsetZ2 = b.LocalOffsetZ
			b.Yaw2 = b.Yaw
			b.Tilt2 = b.Tilt
			b.Pitch2 = 0
		} else if b.ItemId > 2010 && b.ItemId < 2020 {
			b.Tilt = float32(utils.ReadInt32(r))
			b.Pitch = float32(utils.ReadInt32(r))
			b.LocalOffsetX2 = float32(utils.ReadInt32(r))
			b.LocalOffsetY2 = float32(utils.ReadInt32(r))
			b.LocalOffsetZ2 = float32(utils.ReadInt32(r))
			b.Yaw2 = float32(utils.ReadInt32(r))
			b.Tilt2 = float32(utils.ReadInt32(r))
			b.Pitch2 = float32(utils.ReadInt32(r))
		} else {
			b.LocalOffsetX2 = b.LocalOffsetX
			b.LocalOffsetY2 = b.LocalOffsetY
			b.LocalOffsetZ2 = b.LocalOffsetZ
			b.Yaw2 = b.Yaw
			b.Tilt2 = 0
			b.Pitch2 = 0
		}
		b.TempOutputObjIdx = utils.ReadInt32(r)
		b.TempInputObjIdx = utils.ReadInt32(r)
		b.OutputToSlot = utils.ReadInt8(r)
		b.InputFromSlot = utils.ReadInt8(r)
		b.OutputFromSlot = utils.ReadInt8(r)
		b.InputToSlot = utils.ReadInt8(r)
		b.OutputOffset = utils.ReadInt8(r)
		b.InputOffset = utils.ReadInt8(r)
		b.RecipeId = utils.ReadInt16(r)
		b.FilterId = utils.ReadInt16(r)
		l := utils.ReadInt16(r)
		b.Parameters = make([]int32, l)
		for n := 0; n < int(l); n++ {
			b.Parameters[n] = utils.ReadInt32(r)
		}
		b.Content = utils.ReadString(r)
	} else if version <= -101 {
		b.Index = utils.ReadInt32(r)
		b.ItemId = utils.ReadInt16(r)
		b.ModelIndex = utils.ReadInt16(r)
		b.AreaIndex = utils.ReadInt8(r)
		b.LocalOffsetX = float32(utils.ReadInt32(r))
		b.LocalOffsetY = float32(utils.ReadInt32(r))
		b.LocalOffsetZ = float32(utils.ReadInt32(r))
		b.Yaw = float32(utils.ReadInt32(r))
		if b.ItemId > 2000 && b.ItemId < 2010 {
			b.Tilt = float32(utils.ReadInt32(r))
			b.Pitch = 0
			b.LocalOffsetX2 = b.LocalOffsetX
			b.LocalOffsetY2 = b.LocalOffsetY
			b.LocalOffsetZ2 = b.LocalOffsetZ
			b.Yaw2 = b.Yaw
			b.Tilt2 = b.Tilt
			b.Pitch2 = 0
		} else if b.ItemId > 2010 && b.ItemId < 2020 {
			b.Tilt = float32(utils.ReadInt32(r))
			b.Pitch = float32(utils.ReadInt32(r))
			b.LocalOffsetX2 = float32(utils.ReadInt32(r))
			b.LocalOffsetY2 = float32(utils.ReadInt32(r))
			b.LocalOffsetZ2 = float32(utils.ReadInt32(r))
			b.Yaw2 = float32(utils.ReadInt32(r))
			b.Tilt2 = float32(utils.ReadInt32(r))
			b.Pitch2 = float32(utils.ReadInt32(r))
		} else {
			b.LocalOffsetX2 = b.LocalOffsetX
			b.LocalOffsetY2 = b.LocalOffsetY
			b.LocalOffsetZ2 = b.LocalOffsetZ
			b.Yaw2 = b.Yaw
			b.Tilt2 = 0
			b.Pitch2 = 0
		}
		b.TempOutputObjIdx = utils.ReadInt32(r)
		b.TempInputObjIdx = utils.ReadInt32(r)
		b.OutputToSlot = utils.ReadInt8(r)
		b.InputFromSlot = utils.ReadInt8(r)
		b.OutputFromSlot = utils.ReadInt8(r)
		b.InputToSlot = utils.ReadInt8(r)
		b.OutputOffset = utils.ReadInt8(r)
		b.InputOffset = utils.ReadInt8(r)
		b.RecipeId = utils.ReadInt16(r)
		b.FilterId = utils.ReadInt16(r)
		l := utils.ReadInt16(r)
		b.Parameters = make([]int32, l)
		for n := 0; n < int(l); n++ {
			b.Parameters[n] = utils.ReadInt32(r)
		}
	} else if version <= -100 {
		b.Index = utils.ReadInt32(r)
		b.AreaIndex = utils.ReadInt8(r)
		b.LocalOffsetX = float32(utils.ReadInt32(r))
		b.LocalOffsetY = float32(utils.ReadInt32(r))
		b.LocalOffsetZ = float32(utils.ReadInt32(r))
		b.LocalOffsetX2 = float32(utils.ReadInt32(r))
		b.LocalOffsetY2 = float32(utils.ReadInt32(r))
		b.LocalOffsetZ2 = float32(utils.ReadInt32(r))
		b.Pitch = 0
		b.Pitch2 = 0
		b.Yaw = float32(utils.ReadInt32(r))
		b.Yaw2 = float32(utils.ReadInt32(r))
		b.Tilt = float32(utils.ReadInt32(r))
		b.Tilt2 = 0
		b.ItemId = utils.ReadInt16(r)
		b.ModelIndex = utils.ReadInt16(r)
		b.TempOutputObjIdx = utils.ReadInt32(r)
		b.TempInputObjIdx = utils.ReadInt32(r)
		b.OutputToSlot = utils.ReadInt8(r)
		b.InputFromSlot = utils.ReadInt8(r)
		b.OutputFromSlot = utils.ReadInt8(r)
		b.InputToSlot = utils.ReadInt8(r)
		b.OutputOffset = utils.ReadInt8(r)
		b.InputOffset = utils.ReadInt8(r)
		b.RecipeId = utils.ReadInt16(r)
		b.FilterId = utils.ReadInt16(r)
		l := utils.ReadInt16(r)
		b.Parameters = make([]int32, l)
		for n := 0; n < int(l); n++ {
			b.Parameters[n] = utils.ReadInt32(r)
		}
	} else {
		b.Index = utils.ReadInt32(r)
		b.AreaIndex = utils.ReadInt8(r)
		b.LocalOffsetX = float32(utils.ReadInt32(r))
		b.LocalOffsetY = float32(utils.ReadInt32(r))
		b.LocalOffsetZ = float32(utils.ReadInt32(r))
		b.LocalOffsetX2 = float32(utils.ReadInt32(r))
		b.LocalOffsetY2 = float32(utils.ReadInt32(r))
		b.LocalOffsetZ2 = float32(utils.ReadInt32(r))
		b.Pitch = 0
		b.Pitch2 = 0
		b.Yaw = float32(utils.ReadInt32(r))
		b.Yaw2 = float32(utils.ReadInt32(r))
		b.Tilt = 0
		b.Tilt2 = 0
		b.ItemId = utils.ReadInt16(r)
		b.ModelIndex = utils.ReadInt16(r)
		b.TempOutputObjIdx = utils.ReadInt32(r)
		b.TempInputObjIdx = utils.ReadInt32(r)
		b.OutputToSlot = utils.ReadInt8(r)
		b.InputFromSlot = utils.ReadInt8(r)
		b.OutputFromSlot = utils.ReadInt8(r)
		b.InputToSlot = utils.ReadInt8(r)
		b.OutputOffset = utils.ReadInt8(r)
		b.InputOffset = utils.ReadInt8(r)
		b.RecipeId = utils.ReadInt16(r)
		b.FilterId = utils.ReadInt16(r)
		l := utils.ReadInt16(r)
		b.Parameters = make([]int32, l)
		for n := 0; n < int(l); n++ {
			b.Parameters[n] = utils.ReadInt32(r)
		}
	}
	return nil
}

func (b *Building) Write(w io.Writer) error {
	//w.Write(-101);
	//w.Write(this.index);
	//w.Write(this.itemId);
	//w.Write(this.modelIndex);
	//w.Write((sbyte)this.areaIndex);
	//w.Write(this.localOffset_x);
	//w.Write(this.localOffset_y);
	//w.Write(this.localOffset_z);
	//w.Write(this.yaw);
	//if (this.itemId > 2000 && this.itemId < 2010)
	//{
	//	w.Write(this.tilt);
	//}
	//else if (this.itemId > 2010 && this.itemId < 2020)
	//{
	//	w.Write(this.tilt);
	//	w.Write(this.pitch);
	//	w.Write(this.localOffset_x2);
	//	w.Write(this.localOffset_y2);
	//	w.Write(this.localOffset_z2);
	//	w.Write(this.yaw2);
	//	w.Write(this.tilt2);
	//	w.Write(this.pitch2);
	//}
	//w.Write((this.outputObj == null) ? -1 : this.outputObj.index);
	//w.Write((this.inputObj == null) ? -1 : this.inputObj.index);
	//w.Write((sbyte)this.outputToSlot);
	//w.Write((sbyte)this.inputFromSlot);
	//w.Write((sbyte)this.outputFromSlot);
	//w.Write((sbyte)this.inputToSlot);
	//w.Write((sbyte)this.outputOffset);
	//w.Write((sbyte)this.inputOffset);
	//w.Write((short)this.recipeId);
	//w.Write((short)this.filterId);
	//int num = (this.parameters != null) ? this.parameters.Length : 0;
	//w.Write((short)num);
	//for (int i = 0; i < num; i++)
	//{
	//w.Write(this.parameters[i]);
	//}

	utils.WriteInt32(w, -101)
	utils.WriteInt32(w, b.Index)
	utils.WriteInt16(w, b.ItemId)
	utils.WriteInt16(w, b.ModelIndex)
	utils.WriteInt8(w, b.AreaIndex)
	utils.WriteInt32(w, int32(b.LocalOffsetX))
	utils.WriteInt32(w, int32(b.LocalOffsetY))
	utils.WriteInt32(w, int32(b.LocalOffsetZ))
	utils.WriteInt32(w, int32(b.Yaw))
	if b.ItemId > 2000 && b.ItemId < 2010 {
		utils.WriteInt32(w, int32(b.Tilt))
	} else if b.ItemId > 2010 && b.ItemId < 2020 {
		utils.WriteInt32(w, int32(b.Tilt))
		utils.WriteInt32(w, int32(b.Pitch))
		utils.WriteInt32(w, int32(b.LocalOffsetX2))
		utils.WriteInt32(w, int32(b.LocalOffsetY2))
		utils.WriteInt32(w, int32(b.LocalOffsetZ2))
		utils.WriteInt32(w, int32(b.Yaw2))
		utils.WriteInt32(w, int32(b.Tilt2))
		utils.WriteInt32(w, int32(b.Pitch2))
	}
	utils.WriteInt32(w, b.TempOutputObjIdx)
	utils.WriteInt32(w, b.TempInputObjIdx)
	utils.WriteInt8(w, b.OutputToSlot)
	utils.WriteInt8(w, b.InputFromSlot)
	utils.WriteInt8(w, b.OutputFromSlot)
	utils.WriteInt8(w, b.InputToSlot)
	utils.WriteInt8(w, b.OutputOffset)
	utils.WriteInt8(w, b.InputOffset)
	utils.WriteInt16(w, b.RecipeId)
	utils.WriteInt16(w, b.FilterId)
	utils.WriteInt16(w, int16(len(b.Parameters)))
	for _, p := range b.Parameters {
		utils.WriteInt32(w, p)
	}
	return nil
}
