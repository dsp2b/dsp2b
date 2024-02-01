package blueprint

import (
	"errors"
	"fmt"
	"net/url"
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

func Rename(data, name string) (string, error) {
	parse := strings.TrimPrefix(data, "BLUEPRINT:0,")
	if parse == data {
		return "", errors.New("not a blueprint")
	}
	arr := strings.Split(parse, ",")
	if len(arr) != 11 {
		return "", errors.New("not a blueprint")
	}
	arr[9] = url.QueryEscape(name)

	data = "BLUEPRINT:0," + strings.Join(arr, ",")

	arr = strings.Split(data, "\"")
	data = strings.Join(arr[:len(arr)-1], "\"")

	data = data + "\"" + fmt.Sprintf("%X", md5f.MD5Hash(data))

	return data, nil
}
