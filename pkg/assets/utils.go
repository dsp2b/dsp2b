package assets

import (
	"encoding/binary"
	"io"
	"reflect"
)

func ReadStruct(r io.Reader, data any) error {
	ref := reflect.ValueOf(data).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(ReadString(r))
		case reflect.Int32:
			field.SetInt(int64(ReadInt32(r)))
		case reflect.Int64:
			field.SetInt(ReadInt64(r))
		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.Int32:
				field.Set(reflect.ValueOf(ReadInt32Array(r)))
			}
		case reflect.Float32:
			field.SetFloat(float64(ReadInt32(r)))
		case reflect.Bool:
			field.SetBool(ReadBool(r))
		case reflect.Struct:
			ReadStruct(r, field.Addr().Interface())
		}
	}

	return nil
}

func ReadInt32(r io.Reader) int32 {
	var i int32
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
