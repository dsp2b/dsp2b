package blueprint

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dsp2b/dsp2b-go/pkg/blueprint/md5f"

	"github.com/dsp2b/dsp2b-go/pkg/utils"
)

type Blueprint struct {
	Header
	BuildingHeader
	Areas     []Area
	Buildings []Building
	Hash      string
}

func (b *Blueprint) Encode() (string, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("BLUEPRINT:0,%d,%d,%d,%d,%d,%d,%d,%d,%s,%s,%s",
		b.Layout, b.Icon0, b.Icon1, b.Icon2, b.Icon3, b.Icon4, b.Icon5, (b.Time.Unix()+62135596800)*10000*1000,
		b.GameVersion, url.PathEscape(b.ShortDesc), url.PathEscape(b.Desc)))

	buf.WriteString("\"")

	gzipData := bytes.NewBuffer(nil)
	err := func() error {
		gzipW := gzip.NewWriter(gzipData)
		defer gzipW.Close()
		if err := utils.WriteStruct(gzipW, &b.BuildingHeader); err != nil {
			return err
		}

		utils.WriteInt8(gzipW, int8(len(b.Areas)))
		for i := range b.Areas {
			if err := utils.WriteStruct(gzipW, &b.Areas[i]); err != nil {
				return err
			}
		}

		utils.WriteInt32(gzipW, int32(len(b.Buildings)))
		for i := range b.Buildings {
			if err := b.Buildings[i].Write(gzipW); err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return "", err
	}

	buf.WriteString(base64.StdEncoding.EncodeToString(gzipData.Bytes()))

	hash := md5f.MD5Hash(buf.String())

	buf.WriteString("\"")

	buf.WriteString(fmt.Sprintf("%X", hash))

	return buf.String(), nil
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
	hash := s[2]
	removeHash := "BLUEPRINT:0," + strings.Join(arr[:10], ",") + "," + strings.Join(s[:2], "\"")
	// hash校验
	if fmt.Sprintf("%X", md5f.MD5Hash(removeHash)) != hash {
		return errors.New("blueprint hash not match")
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
		// 读取第一个判断版本号
		if err := b.Buildings[i].Parse(r); err != nil {
			return err
		}
	}
	b.Hash = s[2]
	// 蓝图数据
	return nil
}
