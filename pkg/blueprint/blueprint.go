package blueprint

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"github.com/dsp2b/dsp2b-go/pkg/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Blueprint struct {
	Header
	BuildingHeader
	Areas     []Area
	Buildings []Building
	Hash      string
}

// Decode BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,"H4sIAAAAAAAAC2NkQAWMUAxh/2dgOAFlMsKFEWoPSG7Dxv7HYcfwHwpQTWZgAAB4dngncAAAAA=="2881F7A76BAF3A19C17C948A5C773D72
func (b *Blueprint) Decode(data string) error {
	parse := strings.TrimPrefix(data, "BLUEPRINT:0,")
	if parse == data {
		return errors.New("not a blueprint")
	}
	arr := strings.Split(parse, ",")
	if len(arr) != 11 {
		return errors.New("not a blueprint")
	}
	var err error
	s := strings.Split(arr[10], "\"")
	if len(s) != 3 {
		return errors.New("not a blueprint header")
	}
	building := s[1]
	buildingByte := make([]byte, base64.StdEncoding.DecodedLen(len(building)))
	n, err := base64.StdEncoding.Decode(buildingByte, []byte(building))
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(bytes.NewReader(buildingByte[:n]))
	if err != nil {
		return err
	}
	err = b.Header.Decode(arr, s)
	if err != nil {
		return err
	}
	if err := utils.ReadStruct(r, &b.BuildingHeader); err != nil {
		return err
	}
	l := utils.ReadUint8(r)
	b.Areas = make([]Area, l)
	for i := 0; i < int(l); i++ {
		if err := utils.ReadStruct(r, &b.Areas[i]); err != nil {
			return err
		}
	}
	ll := utils.ReadInt32(r)
	b.Buildings = make([]Building, ll)
	for i := 0; i < int(ll); i++ {
		if err := utils.ReadStruct(r, &b.Buildings[i]); err != nil {
			return err
		}
		l := utils.ReadInt16(r)
		b.Buildings[i].Parameters = make([]int32, l)
		for n := 0; n < int(l); n++ {
			b.Buildings[i].Parameters[n] = utils.ReadInt32(r)
		}
	}
	b.Hash = s[2]
	// 蓝图数据
	return nil
}

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

type BuildingHeader struct {
	Version          int32
	CursorOffsetX    int32
	CursorOffsetY    int32
	CursorTargetArea int32
	DragBoxSizeX     int32
	DragBoxSizeY     int32
	PrimaryAreaIdx   int32
}

type Header struct {
	Layout      int
	Icon0       int
	Icon1       int
	Icon2       int
	Icon3       int
	Icon4       int
	Icon5       int
	Time        time.Time
	GameVersion string
	Name        string // shortDesc
	Desc        string // desc
}

func (h *Header) Decode(arr []string, data []string) error {
	if len(arr) != 11 {
		return errors.New("not a blueprint header")
	}
	var err error
	h.Layout, err = strconv.Atoi(arr[0])
	if err != nil {
		return err
	}
	h.Icon0, err = strconv.Atoi(arr[1])
	if err != nil {
		return err
	}
	h.Icon1, err = strconv.Atoi(arr[2])
	if err != nil {
		return err
	}
	h.Icon2, err = strconv.Atoi(arr[3])
	if err != nil {
		return err
	}
	h.Icon3, err = strconv.Atoi(arr[4])
	if err != nil {
		return err
	}
	h.Icon4, err = strconv.Atoi(arr[5])
	if err != nil {
		return err
	}
	h.Icon5, err = strconv.Atoi(arr[6])
	if err != nil {
		return err
	}
	t, err := strconv.ParseInt(arr[7], 10, 64)
	if err != nil {
		return err
	}
	// 1704475921 63839580922
	t = (t / 10000 / 1000) - 62135596800
	h.Time = time.Unix(t, 0)
	h.GameVersion = arr[8]
	h.Name, err = url.QueryUnescape(arr[9])
	if err != nil {
		return err
	}
	h.Desc, err = url.QueryUnescape(data[0])
	if err != nil {
		return err
	}
	return nil
}
