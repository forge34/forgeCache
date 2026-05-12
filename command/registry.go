// Package command
package command

import (
	"strconv"

	"github.com/forge34/forgeCache/resp"
)

type (
	Handler    func(*Store, []resp.Value) resp.Value
	HandlerMap map[string]Handler
)

type Commander struct {
	handlers HandlerMap
	store    *Store
}

func NewCommander() *Commander {
	s := NewStore()
	return &Commander{
		handlers: HandlerMap{
			"PING":   ping,
			"SET":    set,
			"DEL":    deleteKey,
			"GET":    get,
			"INCR":   increase,
			"EXISTS": exists,
			"APPEND": appendStr,
			"DECR":   decrease,
		},
		store: s,
	}
}

func set(s *Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str
	value := args[1].Str

	s.Set(key, &MapValue{Value: value, ExpiresAt: 0})
	return resp.Value{Typ: resp.STRING, Str: "OK"}
}

func get(s *Store, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	key := args[0].Str
	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.BULKSTR, IsNil: true}
	}

	switch v := entry.Value.(type) {
	case string:
		return resp.Value{Typ: resp.BULKSTR, Str: v}
	case int64:
		return resp.Value{Typ: resp.BULKSTR, Str: strconv.FormatInt(v, 10)}
	default:
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}
}

func deleteKey(s *Store, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	deletes := 0
	for _, v := range args {
		ok := s.Delete(v.Str)

		if ok {
			deletes += 1
		}
	}

	return resp.Value{Typ: resp.INTEGER, Num: int64(deletes)}
}

func exists(s *Store, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	count := 0
	for _, v := range args {
		_, ok := s.Get(v.Str)

		if ok {
			count += 1
		}
	}

	return resp.Value{Typ: resp.INTEGER, Num: int64(count)}
}

func increase(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	key := args[0].Str

	v, ok := s.Get(key)
	var current int64

	if !ok {
		current = 0
	} else {
		switch val := v.Value.(type) {
		case int:
			current = int64(val)
		case int64:
			current = val
		case string:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return resp.Value{Typ: resp.ERROR, Str: "ERR value is not an integer"}
			}
			current = n
		default:
			return resp.Value{Typ: resp.ERROR, Str: "ERR wrong type"}
		}
	}

	current++

	s.Set(key, &MapValue{Value: current})

	return resp.Value{
		Typ: resp.INTEGER,
		Num: current,
	}
}

func decrease(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	key := args[0].Str

	v, ok := s.Get(key)
	var current int64

	if !ok {
		current = 0
	} else {
		switch val := v.Value.(type) {
		case int:
			current = int64(val)
		case int64:
			current = val
		case string:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return resp.Value{Typ: resp.ERROR, Str: "ERR value is not an integer"}
			}
			current = n
		default:
			return resp.Value{Typ: resp.ERROR, Str: "ERR wrong type"}
		}
	}

	current--

	s.Set(key, &MapValue{Value: current})

	return resp.Value{
		Typ: resp.INTEGER,
		Num: current,
	}
}

func appendStr(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func ping(s *Store, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: resp.STRING, Str: "PONG"}
	}
	return resp.Value{Typ: resp.BULKSTR, Str: args[0].Str}
}
