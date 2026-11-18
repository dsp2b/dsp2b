package blueprint

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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
	// Tail 0.10.34+ 二进制负载尾部(建筑数据之后: Patch int32 + reformFlag byte + 可选地改数据),
	// 解析时原样保留, 重新编码时写回, 避免丢失地改(reform)蓝图数据
	Tail []byte
}

func (b *Blueprint) Encode() (string, error) {
	buf := bytes.NewBuffer(nil)
	if b.Patch >= 1 {
		// 0.10.34+ 新格式: 新增 author/customVersion/externalFields 三个字段
		buf.WriteString(fmt.Sprintf("BLUEPRINT:1,%d,%d,%d,%d,%d,%d,%d,%d,%s,%s,%s,%s,%s,%s",
			b.Layout, b.Icon0, b.Icon1, b.Icon2, b.Icon3, b.Icon4, b.Icon5, (b.Time.Unix()+62135596800)*10000*1000,
			b.GameVersion, escape(b.ShortDesc), escape(b.Author), escape(b.CustomVersion), escape(b.ExternalFields), escape(b.Desc)))
	} else {
		buf.WriteString(fmt.Sprintf("BLUEPRINT:0,%d,%d,%d,%d,%d,%d,%d,%d,%s,%s,%s",
			b.Layout, b.Icon0, b.Icon1, b.Icon2, b.Icon3, b.Icon4, b.Icon5, (b.Time.Unix()+62135596800)*10000*1000,
			b.GameVersion, escape(b.ShortDesc), escape(b.Desc)))
	}

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
		// 0.10.34+ 二进制尾部: 原样写回; 数据版本>=2 但无尾部时补默认值(Patch=1, 无地改数据)
		if len(b.Tail) > 0 {
			if _, err := gzipW.Write(b.Tail); err != nil {
				return err
			}
		} else if b.BuildingHeader.Version >= 2 {
			utils.WriteInt32(gzipW, 1)
			utils.WriteInt8(gzipW, 0)
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

// Decode BLUEPRINT:0,10,0,0,0,0,0,0,638395809227381327,0.10.28.21014,%E6%96%B0%E8%93%9D%E5%9B%BE,"H4sIA...=="2881F...
// Decode BLUEPRINT:1,10,0,0,0,0,0,0,639305347823681518,0.10.35.29151,%E6%96%B0%E8%93%9D%E5%9B%BE,,,,,,"H4sIA...=="183B...
// (新格式 0.10.34+: 名称后多 author/customVersion/externalFields/desc 四个字段)
func (b *Blueprint) Decode(data string) error {
	if !strings.HasPrefix(data, "BLUEPRINT:") {
		return errors.New("not a blueprint")
	}
	// 头部到第一个引号为止, 与游戏 BlueprintData.HeaderFromBase64String 一致
	quote := strings.IndexByte(data, '"')
	if quote < 0 {
		return errors.New("not a blueprint")
	}
	fields := strings.Split(data[10:quote], ",")
	// hash 为最后一个引号后的 32 个字符
	lastQuote := strings.LastIndexByte(data, '"')
	if lastQuote <= quote || len(data)-lastQuote-1 != 32 {
		return errors.New("not a blueprint")
	}
	hash := data[lastQuote+1:]
	// hash校验
	if fmt.Sprintf("%X", md5f.MD5Hash(data[:lastQuote])) != hash {
		return errors.New("blueprint hash not match")
	}
	building := data[quote+1 : lastQuote]
	buildingByte := make([]byte, base64.StdEncoding.DecodedLen(len(building)))
	n, err := base64.StdEncoding.Decode(buildingByte, []byte(building))
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(bytes.NewReader(buildingByte[:n]))
	if err != nil {
		return err
	}
	err = b.Header.Decode(fields)
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
	// 0.10.34+ 建筑数据之后还有尾部(Patch/地改数据), 原样保留
	if b.Tail, err = io.ReadAll(r); err != nil {
		return err
	}
	b.Hash = hash
	// 蓝图数据
	return nil
}
