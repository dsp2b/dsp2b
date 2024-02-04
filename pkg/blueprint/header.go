package blueprint

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

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
	ShortDesc   string // shortDesc
	Desc        string // desc
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
	h.ShortDesc, err = url.QueryUnescape(arr[9])
	if err != nil {
		return err
	}
	h.Desc, err = url.QueryUnescape(data[0])
	if err != nil {
		return err
	}
	return nil
}
