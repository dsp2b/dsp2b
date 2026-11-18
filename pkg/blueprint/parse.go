package blueprint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dsp2b/dsp2b-go/pkg/blueprint/md5f"
)

func Decode(data string) (Blueprint, error) {
	b := Blueprint{}
	if err := b.Decode(data); err != nil {
		return Blueprint{}, err
	}
	return b, nil
}

// Rename 修改蓝图名称(shortDesc), 保留原格式版本(0/1)与内容, 重新计算hash。
// 名称位置: 头部字段[10](与游戏一致, 旧版解析按 BLUEPRINT:0, 前缀去掉后为 arr[9])
func Rename(data, name string) (string, error) {
	if !strings.HasPrefix(data, "BLUEPRINT:") {
		return "", errors.New("not a blueprint")
	}
	quote := strings.IndexByte(data, '"')
	lastQuote := strings.LastIndexByte(data, '"')
	if quote < 0 || lastQuote <= quote {
		return "", errors.New("not a blueprint")
	}
	fields := strings.Split(data[10:quote], ",")
	if len(fields) < 12 {
		return "", errors.New("not a blueprint")
	}
	fields[10] = escape(name)

	noHash := "BLUEPRINT:" + strings.Join(fields, ",") + "\"" + data[quote+1:lastQuote]

	return noHash + "\"" + fmt.Sprintf("%X", md5f.MD5Hash(noHash)), nil
}
