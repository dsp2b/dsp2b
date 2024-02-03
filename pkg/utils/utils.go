package utils

import (
	"encoding/binary"
	"io"
	"reflect"
	"strings"
)

func ReadStruct(r io.Reader, data any) error {
	ref := reflect.ValueOf(data).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := ref.Type().Field(i).Tag.Get("binary")
		if tag == "-" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(ReadString(r))
		case reflect.Int8:
			field.SetInt(int64(ReadInt8(r)))
		case reflect.Int16:
			field.SetInt(int64(ReadInt16(r)))
		case reflect.Int32:
			field.SetInt(int64(ReadInt32(r)))
		case reflect.Int64:
			field.SetInt(ReadInt64(r))
		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.Int32:
				field.Set(reflect.ValueOf(ReadInt32Array(r)))
			case reflect.Float64:
				field.Set(reflect.ValueOf(ReadFloat64Array(r)))
			}
		case reflect.Float32:
			field.SetFloat(float64(ReadInt32(r)))
		case reflect.Bool:
			field.SetBool(ReadBool(r))
		case reflect.Struct:
			if err := ReadStruct(r, field.Addr().Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadInt32(r io.Reader) int32 {
	var i int32
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadUint8(r io.Reader) uint8 {
	var i uint8
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadInt8(r io.Reader) int8 {
	var i int8
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadInt16(r io.Reader) int16 {
	var i int16
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadInt64(r io.Reader) int64 {
	var i int64
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadInt32Array(r io.Reader) []int32 {
	var n int32
	_ = binary.Read(r, binary.LittleEndian, &n)
	ret := make([]int32, n)
	for i := 0; i < int(n); i++ {
		_ = binary.Read(r, binary.LittleEndian, &ret[i])
	}
	return ret
}

func ReadFloat64Array(r io.Reader) []float64 {
	var n int32
	_ = binary.Read(r, binary.LittleEndian, &n)
	ret := make([]float64, n)
	for i := 0; i < int(n); i++ {
		_ = binary.Read(r, binary.LittleEndian, &ret[i])
	}
	return ret
}

func ReadString(r io.Reader) string {
	var l int32
	_ = binary.Read(r, binary.LittleEndian, &l)
	if l == 0 {
		return ""
	}
	b := make([]byte, l)
	_ = binary.Read(r, binary.LittleEndian, &b)
	//对其
	empty := make([]byte, -l&3)
	_ = binary.Read(r, binary.LittleEndian, &empty)
	return string(b)
}

func ReadBool(r io.Reader) bool {
	var b int32
	_ = binary.Read(r, binary.LittleEndian, &b)
	return b == 1
}

func ToPathUnderline(s string) string {
	ss := strings.Split(s, "/")
	for i, s := range ss {
		ss[i] = ToUnderline(s)
	}
	return strings.Join(ss, "/")
}

func ToUnderline(s string) string {
	var ret []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			ret = append(ret, '_')
		}
		ret = append(ret, r)
	}
	return string(ret)
}
