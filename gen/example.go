package gen

import "fjson/codec"

type ExampleNested struct {
	Id       int32
	Int      int
	Str      string
	IntArray []int
}

type Example struct {
	Id          int32
	Int         int
	Str         string
	IntArray    []int
	StrArray    []string
	Nested      ExampleNested
	NestedArray []ExampleNested
}

func decodeExample(r codec.ByteReader, out *Example) {
	codec.ReadStructBegin(r)
	out.Int, _ = codec.ReadInt(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.Str, _ = codec.ReadString(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.IntArray, _ = codec.ReadSlice(r, codec.ReadInt)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.StrArray, _ = codec.ReadSlice(r, codec.ReadString)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	decodeExampleNest(r, &out.Nested)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.NestedArray, _ = codec.ReadSlice(r, func(r codec.ByteReader) (ExampleNested, error) {
		nest := ExampleNested{}
		decodeExampleNest(r, &nest)
		return nest, nil
	})
	codec.ReadStructEnd(r)
}

func decodeExampleNest(r codec.ByteReader, out *ExampleNested) {
	codec.ReadStructBegin(r)
	out.Int, _ = codec.ReadInt(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.Str, _ = codec.ReadString(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.IntArray, _ = codec.ReadSlice(r, codec.ReadInt)
	codec.ReadStructEnd(r)
}
