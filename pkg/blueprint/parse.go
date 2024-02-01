package blueprint

import (
	"errors"
	"net/url"
	"strings"
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
	return "BLUEPRINT:0," + strings.Join(arr, ","), nil
}
