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

func WriteStruct(w io.Writer, data any) error {
	ref := reflect.ValueOf(data).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := ref.Type().Field(i).Tag.Get("binary")
		if tag == "-" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			WriteString(w, field.String())
		case reflect.Int8:
			WriteInt8(w, int8(field.Int()))
		case reflect.Int16:
			WriteInt16(w, int16(field.Int()))
		case reflect.Int32:
			WriteInt32(w, int32(field.Int()))
		case reflect.Int64:
			WriteInt64(w, field.Int())
		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.Int32:
				WriteInt32Array(w, field.Interface().([]int32))
			case reflect.Float64:
				WriteFloat64Array(w, field.Interface().([]float64))
			}
		case reflect.Float32:
			WriteInt32(w, int32(field.Float()))
		case reflect.Bool:
			if field.Bool() {
				WriteInt32(w, 1)
			} else {
				WriteInt32(w, 0)
			}
		case reflect.Struct:
			if err := WriteStruct(w, field.Addr().Interface()); err != nil {
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

func WriteInt32(w io.Writer, i int32) {
	_ = binary.Write(w, binary.LittleEndian, i)
}

func ReadInt8(r io.Reader) int8 {
	var i int8
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func WriteInt8(w io.Writer, i int8) {
	_ = binary.Write(w, binary.LittleEndian, i)
}

func ReadInt16(r io.Reader) int16 {
	var i int16
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func WriteInt16(w io.Writer, i int16) {
	_ = binary.Write(w, binary.LittleEndian, i)
}

func ReadInt64(r io.Reader) int64 {
	var i int64
	_ = binary.Read(r, binary.LittleEndian, &i)
	return i
}

func WriteInt64(w io.Writer, i int64) {
	_ = binary.Write(w, binary.LittleEndian, i)
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

func WriteInt32Array(w io.Writer, arr []int32) {
	_ = binary.Write(w, binary.LittleEndian, int32(len(arr)))
	for _, i := range arr {
		_ = binary.Write(w, binary.LittleEndian, i)
	}
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

func WriteFloat64Array(w io.Writer, arr []float64) {
	_ = binary.Write(w, binary.LittleEndian, int32(len(arr)))
	for _, i := range arr {
		_ = binary.Write(w, binary.LittleEndian, i)
	}
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

func WriteString(w io.Writer, s string) {
	_ = binary.Write(w, binary.LittleEndian, int32(len(s)))
	_ = binary.Write(w, binary.LittleEndian, []byte(s))
	empty := make([]byte, -len(s)&3)
	_ = binary.Write(w, binary.LittleEndian, empty)
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
	ret := make([]rune, 0)
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			ret = append(ret, '_')
		}
		ret = append(ret, r)
	}
	return string(ret)
}
