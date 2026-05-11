// Package resp
package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	Reader *bufio.Reader
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{Reader: bufio.NewReader(rd)}
}

func (p *Parser) Read() (Value, error) {
	typ, err := p.Reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch ValueType(typ) {
	case STRING:
		return p.ReadString()
	case BULKSTR:
		return p.ReadBulk()
	case INTEGER:
		return p.ReadInteger()
	case ERROR:
		return p.ReadError()
	case ARRAY:
		return p.ReadArray()
	default:
		fmt.Printf("Unknown type: %v", string(typ))
		return Value{}, nil
	}
}

func (p *Parser) ReadString() (Value, error) {
	line, err := p.Reader.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}

	return Value{Typ: STRING, Str: string(line[:len(line)-2])}, nil
}

func (p *Parser) ReadInteger() (Value, error) {
	line, err := p.Reader.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}
	num, err := strconv.ParseInt(string(line[:len(line)-2]), 10, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{Typ: INTEGER, Num: num}, nil
}

func (p *Parser) ReadBulk() (Value, error) {
	line, err := p.Reader.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}

	line = line[:len(line)-2]

	n, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, err
	}

	if n == -1 {
		return Value{Typ: BULKSTR, IsNil: true}, nil
	}

	length := int(n)
	buf := make([]byte, length)

	_, err = io.ReadFull(p.Reader, buf)
	if err != nil {
		return Value{}, err
	}

	if _, err := p.Reader.ReadByte(); err != nil {
		return Value{}, err
	}
	if _, err := p.Reader.ReadByte(); err != nil {
		return Value{}, err
	}

	return Value{Typ: BULKSTR, Str: string(buf)}, nil
}

func (p *Parser) ReadArray() (Value, error) {
	line, err := p.Reader.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}

	line = line[:len(line)-2]

	n, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, err
	}

	values := make([]Value, 0, n)

	for i := 0; i < int(n); i++ {
		v, err := p.Read()
		if err != nil {
			return Value{}, nil
		}

		values = append(values, v)
	}
	return Value{Typ: ARRAY, Array: values}, nil
}
func (p *Parser) ReadError() (Value, error) { return Value{}, nil }
