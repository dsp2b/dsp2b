package blueprint

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Blueprint struct {
	Header   *Header
	Building *Building
}

// Decode BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,"H4sIAAAAAAAAC2NkQAWMUAxh/2dgOAFlMsKFEWoPSG7Dxv7HYcfwHwpQTWZgAAB4dngncAAAAA=="2881F7A76BAF3A19C17C948A5C773D72
func (b *Blueprint) Decode(data string) (*Blueprint, error) {
	parse := strings.TrimPrefix(data, "BLUEPRINT:0,")
	if parse == data {
		return nil, errors.New("not a blueprint")
	}
	arr := strings.Split(parse, ",")
	if len(arr) != 11 {
		return nil, errors.New("not a blueprint")
	}
	var err error
	header := &Header{}
	b.Header, err = header.Decode(arr[:10])
	if err != nil {
		return nil, err
	}
	building := &Building{}
	b.Building, err = building.Decode(arr[10])
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Building struct {
	Version       int
	CursorOffsetX int
	Building      string
	Hash          string
}

func (b *Building) Decode(data string) (*Building, error) {
	arr := strings.Split(data, "\"")
	if len(arr) != 3 {
		return nil, errors.New("not a building")
	}
	b.Building = arr[0]
	b.Hash = arr[1]
	return b, nil
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

func (h *Header) Decode(arr []string) (*Header, error) {
	if len(arr) != 10 {
		return nil, errors.New("not a blueprint header")
	}
	var err error
	h.Layout, err = strconv.Atoi(arr[0])
	if err != nil {
		return nil, err
	}
	h.Icon0, err = strconv.Atoi(arr[1])
	if err != nil {
		return nil, err
	}
	h.Icon1, err = strconv.Atoi(arr[2])
	if err != nil {
		return nil, err
	}
	h.Icon2, err = strconv.Atoi(arr[3])
	if err != nil {
		return nil, err
	}
	h.Icon3, err = strconv.Atoi(arr[4])
	if err != nil {
		return nil, err
	}
	h.Icon4, err = strconv.Atoi(arr[5])
	if err != nil {
		return nil, err
	}
	h.Icon5, err = strconv.Atoi(arr[6])
	if err != nil {
		return nil, err
	}
	t, err := strconv.ParseInt(arr[7], 10, 64)
	if err != nil {
		return nil, err
	}
	t = t / 10000 / 1000
	h.Time = time.Unix(t, 0)
	h.GameVersion = arr[8]
	h.Name, err = url.QueryUnescape(arr[9])
	if err != nil {
		return nil, err
	}
	h.Desc, err = url.QueryUnescape(arr[10])
	if err != nil {
		return nil, err
	}
	return h, nil
}
