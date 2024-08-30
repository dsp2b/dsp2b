package blueprint

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
	Version        int32 // 版本 =-100时，增加了Tilt字段
	Index          int32
	AreaIndex      int8
	LocalOffsetX   float32
	LocalOffsetY   float32
	LocalOffsetZ   float32
	LocalOffsetX2  float32
	LocalOffsetY2  float32
	LocalOffsetZ2  float32
	Yaw            float32
	Yaw2           float32
	Tilt           float32 `binary:"version=-100"`
	ItemId         int16
	ModelIndex     int16
	OutputObj      int32
	InputObj       int32
	OutputToSlot   int8
	InputFromSlot  int8
	OutputFromSlot int8
	InputToSlot    int8
	OutputOffset   int8
	InputOffset    int8
	RecipeId       int16
	FilterId       int16
	Parameters     []int32 `binary:"-"`
}
