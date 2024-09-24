package codec

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func BenchmarkGob(b *testing.B) {
	var e []*Example
	for i := 0; i < b.N; i++ {
		decoder := gob.NewDecoder(bytes.NewBuffer(gobBytes))
		err := decoder.Decode(&e)
		if err != nil {
			b.Fatal(err)
		}
	}
}
