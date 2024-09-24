package gen

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func ParserStruct(v reflect.Type, w io.Writer) (err error) {
	switch v.Kind() {
	case reflect.Bool:
		_, err = w.Write([]byte(fmt.Sprintf("%t", v.Bool())))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(v.Int()))
		_, err = w.Write(bs)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, v.Uint())
		_, err = w.Write(bs)
	case reflect.Float32, reflect.Float64:
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(v.Float()))
		_, err = w.Write(bs)
	case reflect.String:
		l := uint16(len(v.String()))
		if l > 0 {
			bs := make([]byte, 3)
			bs[0] = '"'
			binary.BigEndian.PutUint16(bs[1:], l)
			_, err = w.Write(bs)
			if err != nil {
				return
			}
			_, err = w.Write(unsafe.Slice(unsafe.StringData(v.String()), len(v.String())))
			if err != nil {
				return
			}
			_, err = w.Write(bs[:1])
			if err != nil {
				return
			}
		} else {
			_, err = w.Write([]byte(`""`))
			if err != nil {
				return
			}
		}
	case reflect.Struct:
		_, err = w.Write([]byte("{"))
		if err != nil {
			return
		}
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				_, err = w.Write([]byte(","))
				if err != nil {
					return
				}
			}
			err = GenValue(v.Field(i), w)
			if err != nil {
				return
			}
		}
		_, err = w.Write([]byte("}"))
		if err != nil {
			return
		}
	case reflect.Slice, reflect.Array:
		l := uint16(v.Len())
		if l > 0 {
			bs := make([]byte, 3)
			bs[0] = '['
			binary.BigEndian.PutUint16(bs, l)
			_, err = w.Write(bs)
			if err != nil {
				return
			}
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					_, err = w.Write([]byte(","))
				}
				err = GenValue(v.Index(i), w)
				if err != nil {
					return
				}
			}
			_, err = w.Write(bs[:1])
			if err != nil {
				return
			}
		} else {
			_, err = w.Write([]byte("[]"))
		}
	case reflect.Pointer:
		err = GenValue(v.Elem(), w)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unsupported type: %s, v:%v", v.Kind(), v)
	}
	return
}
