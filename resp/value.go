package resp

import "fmt"

type ValueType byte

const (
	STRING  ValueType = '+'
	ERROR   ValueType = '-'
	INTEGER ValueType = ':'
	BULKSTR ValueType = '$'
	ARRAY   ValueType = '*'
)

type Value struct {
	Typ   ValueType
	Str   string
	Num   int64
	Array []Value
}

func (v Value) Pretty(indent string) string {
	switch v.Typ {
	case ARRAY:
		s := "ARRAY[\n"
		for _, item := range v.Array {
			s += indent + "  " + item.Pretty(indent+"  ") + "\n"
		}
		s += indent + "]"
		return s

	case BULKSTR:
		return fmt.Sprintf("BULK(%q)", v.Str)

	case STRING:
		return fmt.Sprintf("STRING(%q)", v.Str)

	case INTEGER:
		return fmt.Sprintf("INT(%d)", v.Num)

	default:
		return "UNKNOWN"
	}
}
