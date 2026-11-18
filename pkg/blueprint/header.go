package blueprint

import (
	"errors"
	"strconv"
	"time"
)

type Header struct {
	Patch          int // 头部格式版本: 0(<=0.10.33旧格式) / 1(0.10.34+新格式)
	Layout         int
	Icon0          int
	Icon1          int
	Icon2          int
	Icon3          int
	Icon4          int
	Icon5          int
	Time           time.Time
	GameVersion    string
	ShortDesc      string // shortDesc
	Desc           string // desc
	Author         string // 0.10.34+ 蓝图作者
	CustomVersion  string // 0.10.34+ 自定义版本号
	ExternalFields string // 0.10.34+ 自定义属性项
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

// Decode 解析头部字段, fields 为按逗号切分的完整头部(含 fields[0]=格式版本),
// 与游戏 BlueprintData.HeaderFromBase64String 的字段索引一致:
//
//	[0]=patch [1]=layout [2..6]=icon0..icon4 [7]=保留0 [8]=time(.NET Ticks)
//	[9]=gameVersion [10]=shortDesc
//	patch>=1: [11]=author [12]=customVersion [13]=externalFields [14]=desc
//	patch=0:  [11]=desc
func (h *Header) Decode(fields []string) error {
	if len(fields) < 12 {
		return errors.New("not a blueprint header")
	}
	var err error
	h.Patch, err = strconv.Atoi(fields[0])
	if err != nil {
		return err
	}
	if h.Patch < 0 || h.Patch > 1 {
		return errors.New("not a blueprint header")
	}
	if h.Patch >= 1 && len(fields) < 15 {
		return errors.New("not a blueprint header")
	}
	h.Layout, err = strconv.Atoi(fields[1])
	if err != nil {
		return err
	}
	h.Icon0, err = strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	h.Icon1, err = strconv.Atoi(fields[3])
	if err != nil {
		return err
	}
	h.Icon2, err = strconv.Atoi(fields[4])
	if err != nil {
		return err
	}
	h.Icon3, err = strconv.Atoi(fields[5])
	if err != nil {
		return err
	}
	h.Icon4, err = strconv.Atoi(fields[6])
	if err != nil {
		return err
	}
	h.Icon5, err = strconv.Atoi(fields[7])
	if err != nil {
		return err
	}
	t, err := strconv.ParseInt(fields[8], 10, 64)
	if err != nil {
		return err
	}
	// 1704475921 63839580922 (100ns单位, .NET Ticks, 起0001-01-01)
	t = (t / 10000 / 1000) - 62135596800
	h.Time = time.Unix(t, 0)
	h.GameVersion = fields[9]
	h.ShortDesc = unescape(fields[10])
	if h.Patch >= 1 {
		h.Author = unescape(fields[11])
		h.CustomVersion = unescape(fields[12])
		h.ExternalFields = unescape(fields[13])
		h.Desc = unescape(fields[14])
	} else {
		h.Desc = unescape(fields[11])
	}
	return nil
}
