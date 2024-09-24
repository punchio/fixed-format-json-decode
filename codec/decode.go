package codec

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	errInvalidChar    = errors.New("ReadString invalid char")
	errArrayNotFinish = errors.New("array not finish")
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

func ReadInt(r ByteReader) (int, error) {
	if ReadComma(r) {
		_ = r.UnreadByte()
		return 0, nil
	}
	_, err := r.Read(lenInt[:])
	if err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint64(lenInt[:])), nil
}

var lenSlice [2]byte

func ReadSlice[T any](r ByteReader, f func(ByteReader) (T, error)) ([]T, error) {
	if ReadComma(r) {
		_ = r.UnreadByte()
		return nil, nil
	}
	var arr []T
	ReadArrayBegin(r)
	if !ReadArrayEnd(r) {
		_ = r.UnreadByte()
		_, err := r.Read(lenSlice[:])
		if err != nil {
			return arr, err
		}
		l := binary.BigEndian.Uint16(lenSlice[:])
		arr = make([]T, l)
		for i := 0; i < int(l); i++ {
			arr[i], err = f(r)
			if err != nil {
				return nil, err
			}
			if i != int(l)-1 && !ReadComma(r) {
				return nil, errArrayNotFinish
			}
		}

		if !ReadArrayEnd(r) {
			return nil, errArrayNotFinish
		}
	}

	return arr, nil
}

func readQuote(r ByteReader) bool {
	readByte, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return readByte == '"'
}

var strBuf [1024]byte

func ReadString(r ByteReader) (string, error) {
	if ReadComma(r) {
		_ = r.UnreadByte()
		return "", nil
	}
	if !readQuote(r) {
		return "", errInvalidChar
	}

	if readQuote(r) {
		return "", nil
	}
	_ = r.UnreadByte()

	_, err := r.Read(lenSlice[:])
	if err != nil {
		return "", err
	}
	l := binary.BigEndian.Uint16(lenSlice[:])
	_, err = r.Read(strBuf[:l])
	if err != nil {
		return "", err
	}
	if !readQuote(r) {
		return "", errInvalidChar
	}
	return string(strBuf[:l]), nil
}
