package codec

import "fjson/codec"

func decodeExample(r codec.ByteReader, out *Example) {
	codec.ReadStructBegin(r)
	out.Int = codec.ReadInt(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.Str = codec.ReadString(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.IntArray = codec.ReadSlice(r, codec.ReadInt)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.StrArray = codec.ReadSlice(r, codec.ReadString)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	decodeExampleNest(r, &out.Nested)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.NestedArray = codec.ReadSlice(r, func(r codec.ByteReader) ExampleNested {
		nest := ExampleNested{}
		decodeExampleNest(r, &nest)
		return nest
	})
	codec.ReadStructEnd(r)
}

func decodeExampleNest(r codec.ByteReader, out *ExampleNested) {
	codec.ReadStructBegin(r)
	out.Int = codec.ReadInt(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.Str = codec.ReadString(r)
	if !codec.ReadComma(r) {
		panic("need comma")
	}
	out.IntArray = codec.ReadSlice(r, codec.ReadInt)
	codec.ReadStructEnd(r)
}
