package codec

import (
	"encoding/binary"
	"io"
)

type ByteReader interface {
	io.Seeker
	io.ByteScanner
	io.Reader
}

func ReadComma(r ByteReader) bool {
	readByte, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte == ','
}

func ReadStructBegin(r ByteReader) {
	readByte, err := r.ReadByte()
	if err != nil || readByte != '{' {
		panic(err)
	}
}

func ReadStructEnd(r ByteReader) {
	readByte, err := r.ReadByte()
	if err != nil || readByte != '}' {
		panic(err)
	}
}

func ReadArrayBegin(r ByteReader) {
	readByte, err := r.ReadByte()
	if err != nil || readByte != '[' {
		panic(err)
	}
}

func ReadArrayEnd(r ByteReader) bool {
	readByte, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte == ']'
}

var lenInt [8]byte

func ReadInt(r ByteReader) int {
	_, err := r.Read(lenInt[:])
	if err != nil {
		panic(err)
	}

	return int(binary.BigEndian.Uint64(lenInt[:]))
}

var lenSlice [2]byte

func ReadSlice[T any](r ByteReader, f func(ByteReader) T) []T {
	var arr []T
	ReadArrayBegin(r)
	if !ReadArrayEnd(r) {
		_ = r.UnreadByte()
		_, err := r.Read(lenSlice[:])
		if err != nil {
			panic(err)
		}
		l := binary.BigEndian.Uint16(lenSlice[:])
		arr = make([]T, l)
		for i := 0; i < int(l); i++ {
			arr[i] = f(r)
			if i != int(l)-1 && !ReadComma(r) {
				panic("array not finish")
			}
		}

		if !ReadArrayEnd(r) {
			panic("array not finish")
		}
	}

	return arr
}

func ReadQuote(r ByteReader) bool {
	readByte, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte == '"'
}

var strBuf [1024]byte

func ReadString(r ByteReader) string {
	if !ReadQuote(r) {
		panic("ReadString invalid char")
	}

	if ReadQuote(r) {
		return ""
	}
	_ = r.UnreadByte()

	_, err := r.Read(lenSlice[:])
	if err != nil {
		panic(err)
	}
	l := binary.BigEndian.Uint16(lenSlice[:])
	_, err = r.Read(strBuf[:l])
	if err != nil {
		panic(err)
	}
	if !ReadQuote(r) {
		panic("ReadString invalid char")
	}
	return string(strBuf[:l])
}
