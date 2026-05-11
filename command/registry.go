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
	entry := s.Get(key)

	if entry == nil {
		return resp.Value{Typ: resp.BULKSTR, IsNil: true}
	}

	switch v := entry.Value.(type) {
	case string:
		return resp.Value{Typ: resp.BULKSTR, Str: v}
	case int64:
		return resp.Value{Typ: resp.BULKSTR, Str: strconv.FormatInt(v, 10)}
	case int:
		return resp.Value{Typ: resp.BULKSTR, Str: strconv.Itoa(v)}
	default:
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}
}

func deleteKey(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func exists(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func increase(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func decrease(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func appendStr(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func ping(s *Store, args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}
