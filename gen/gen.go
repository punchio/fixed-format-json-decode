package gen

import (
	"io"
	"regexp"
	"strings"
)

type primitiveType int

const (
	None primitiveType = iota
	Int
	Float
	String
	Bool
	Array
	Map
	Struct
	Pointer
	Import
)

var (
	arrayMatch = regexp.MustCompile(`array<(.+)>`)
	mapMatch   = regexp.MustCompile(`map<(.+),(.+)>`)
)

type FieldInfo struct {
	Name    string `toml:"name"`
	Type    string `toml:"type"`
	Comment string `toml:"comment"`
	Filter  string `toml:"filter"`
}

type StructInfo struct {
	Name    string       `toml:"name"`
	Comment string       `toml:"comment"`
	Filter  string       `toml:"filter"`
	Fields  []*FieldInfo `toml:"field_list"`
}

type ImportInfo struct {
	Name       string `toml:"name"`
	Type       string `toml:"type"`
	ClientType string `toml:"client"`
	ServerType string `toml:"server"`
}

type TypeMgr struct {
	Items     []*StructInfo `toml:"type_list"`
	ItemMap   map[string]*StructInfo
	Imports   []*ImportInfo `toml:"import_list"`
	ImportMap map[string]*ImportInfo
}

func Parse(typeMgr *TypeMgr) error {
	return nil
}

func genType(typeMgr *TypeMgr, st *StructInfo, writer io.Writer) error {
	//writer.Write()
	return nil
}

func genInt(fieldName string, writer io.Writer) error {
	_, err := writer.Write([]byte(`s.` + fieldName + `, err = codec.ReadInt(r)
	if err != nil {
		return err
	}
`))
	return err
}

func genString(fieldName string, writer io.Writer) error {
	_, err := writer.Write([]byte(`s.` + fieldName + `, err = codec.ReadString(r)
	if err != nil {
		return err
	}
`))
	return err
}

const sliceTmpl = `
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
`

func genSlice(fieldName string, writer io.Writer) error {
	_, err := writer.Write([]byte(`s.` + fieldName + `, err = codec.ReadSlice(r, )
	if err != nil {
		return err
	}
`))
	return err
}

func getTypeInfo(typeMgr *TypeMgr, typeName string) (primitiveType, string, string) {
	if primitive := isPrimitive(typeName); primitive != None {
		return primitive, typeName, ""
	} else if subType, ok := isArray(typeName); ok {
		return Array, subType, ""
	} else if keyType, valueType, ok2 := isMap(typeName); ok2 {
		return Map, keyType, valueType
	} else if subType, ok = isPointer(typeName); ok {
		return Pointer, subType, ""
	} else if subType, ok = isImport(typeMgr, typeName); ok {
		return Import, subType, ""
	} else if ok = isStruct(typeMgr, typeName); ok {
		return Struct, typeName, ""
	} else {
		return None, "", ""
	}
}

func isPrimitive(typeName string) primitiveType {
	switch typeName {
	case "int32", "int64", "uint32", "uint64":
		return Int
	case "float":
		return Float
	case "bool":
		return Bool
	case "string":
		return String
	default:
		return None
	}
}

func isArray(typeName string) (string, bool) {
	m := arrayMatch.FindStringSubmatch(typeName)
	if len(m) == 2 {
		return m[1], true
	}
	return "", false
}

func isMap(typeName string) (string, string, bool) {
	m := mapMatch.FindStringSubmatch(typeName)
	if len(m) == 3 {
		return m[1], m[2], true
	}
	return "", "", false
}

func isPointer(typeName string) (string, bool) {
	if strings.HasPrefix(typeName, "*") {
		return strings.TrimLeft(typeName, "*"), true
	}
	return typeName, false
}

func isImport(typeMgr *TypeMgr, typeName string) (string, bool) {
	info, ok := typeMgr.ImportMap[typeName]
	if !ok {
		return "", false
	}
	return info.Type, true
}

func isStruct(typeMgr *TypeMgr, typeName string) bool {
	_, ok := typeMgr.ItemMap[typeName]
	return ok
}

/*
	(t *ExampleNested) Unmarshal(r codec.ByteReader) error {
		codec.ReadStructBegin(r)
		codec.ReadStructEnd(r)
	}
*/
