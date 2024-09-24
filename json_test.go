package codec

import (
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func BenchmarkJson(b *testing.B) {
	var e []*Example
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(jsonBytes, &e)
		if err != nil {
			b.Fatal(err)
		}
	}
}
