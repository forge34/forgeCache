package resp

import (
	"fmt"
	"strconv"
)

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
	IsNil bool
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case STRING:
		return v.MarshalString()
	case INTEGER:
		return v.MarshalInteger()
	case ERROR:
		return v.MarshalError()
	case BULKSTR:
		return v.MarshalBulkStr()
	case ARRAY:
		return v.MarshalArray()
	default:

		return []byte{}
	}
}

func (v Value) MarshalArray() []byte {
	var p []byte
	n := len(v.Array)
	p = append(p, '*')
	p = append(p, []byte(strconv.Itoa(n))...)
	p = append(p, '\r', '\n')

	for i := range n {
		p2 := v.Array[i].Marshal()
		p = append(p, p2...)
	}
	return p
}

func (v Value) MarshalBulkStr() []byte {
	if v.IsNil {
		return []byte("$-1\r\n")
	}
	var p []byte
	n := len(v.Str)
	p = append(p, '$')
	p = append(p, []byte(strconv.Itoa(n))...)
	p = append(p, '\r', '\n')
	p = append(p, []byte(v.Str)...)
	p = append(p, '\r', '\n')
	return p
}

func (v Value) MarshalInteger() []byte {
	str := ":" + strconv.FormatInt(v.Num, 10) + "\r\n"

	return []byte(str)
}

func (v Value) MarshalError() []byte {
	str := "-" + v.Str + "\r\n"

	return []byte(str)
}

func (v Value) MarshalString() []byte {
	str := "+" + v.Str + "\r\n"

	return []byte(str)
}

func (v Value) Pretty(indent string) string {
	switch v.Typ {
	case ARRAY:
		s := "ARRAY[\n"
		for _, item := range v.Array {
			s += indent + "  " + item.Pretty(indent+"  ") + "\n"
		}
		s += indent + "]" + "\n"
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
