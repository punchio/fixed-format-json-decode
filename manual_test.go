package codec

import (
	"bytes"
	"fjson/codec"
	"testing"
)

func BenchmarkManual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(manualBytes)
		examples = codec.ReadSlice[*Example](r, func(reader codec.ByteReader) *Example {
			e := &Example{}
			decodeExample(reader, e)
			return e
		})
	}
}
