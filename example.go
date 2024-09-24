package codec

import (
	"bytes"
	"encoding/gob"
	"fjson/codec"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type ExampleNested struct {
	Int      int
	Str      string
	IntArray []int
}

type Example struct {
	Int         int
	Str         string
	IntArray    []int
	StrArray    []string
	Nested      ExampleNested
	NestedArray []ExampleNested
}

var examples []*Example
var gobBytes = make([]byte, 0, 100000)
var manualBytes = make([]byte, 0, 100000)
var jsonBytes = make([]byte, 0, 100000)
var count = 100

func init() {
	example := Example{}
	example.Int = 12345
	example.Str = `Hello, "world"!`
	example.IntArray = []int{1, 2, 3, 4, 5}
	example.StrArray = []string{"one", "two", "three"}
	example.Nested = ExampleNested{
		Int:      67890,
		Str:      "Nested",
		IntArray: []int{6, 7, 8, 9, 0},
	}
	example.NestedArray = []ExampleNested{
		{Int: 11111, Str: `"Nested1"`},
		{Int: 22222, Str: "Nested2"},
		{Int: 33333, Str: "Nested3"},
	}

	for i := 0; i < count; i++ {
		tmp := example
		examples = append(examples, &tmp)
	}

	buf := bytes.NewBuffer(gobBytes)
	err := gob.NewEncoder(buf).Encode(examples)
	if err != nil {
		panic(err)
	}
	gobBytes = buf.Bytes()

	jsonBytes, err = jsoniter.Marshal(examples)
	if err != nil {
		panic(err)
	}
	buf = bytes.NewBuffer(manualBytes)
	err = codec.GenValue(reflect.ValueOf(examples), buf)
	if err != nil {
		panic(err)
	}
	manualBytes = buf.Bytes()
}
