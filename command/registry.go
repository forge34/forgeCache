// Package command
package command

import (
	"strconv"
	"strings"
	"time"

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
		store: s,
		handlers: HandlerMap{
			"PING":   ping,
			"SET":    set,
			"DEL":    deleteKey,
			"GET":    get,
			"INCR":   increase,
			"EXISTS": exists,
			"APPEND": appendStr,
			"DECR":   decrease,
			"EXPIRE": expire,
			"TTL":    checkTTL,
		},
	}
}

func checkTTL(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{
			Typ: resp.ERROR,
			Str: "ERR wrong number of arguments",
		}
	}

	key := args[0].Str

	v, ok := s.Get(key)
	if !ok {
		return resp.Value{
			Typ: resp.INTEGER,
			Num: -2,
		}
	}

	if v.ExpiresAt == 0 {
		return resp.Value{
			Typ: resp.INTEGER,
			Num: -1,
		}
	}

	timeLeft := v.ExpiresAt - time.Now().UnixNano()

	return resp.Value{
		Typ: resp.INTEGER,
		Num: timeLeft / int64(time.Second),
	}
}

func expire(s *Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str
	value := args[1].Str
	ttl, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return resp.Value{Typ: resp.ERROR, Str: "ERR invalid expire time"}
	}
	ok := s.UpdateExpiry(key, time.Now().Add(time.Duration(ttl)*time.Second))

	if !ok {
		return resp.Value{Typ: resp.INTEGER, Num: 0}
	}

	return resp.Value{Typ: resp.INTEGER, Num: 1}
}

func set(s *Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str
	value := args[1].Str
	var expiresAt int64
	if len(args) >= 4 && strings.ToUpper(args[2].Str) == "EX" {
		ttl, err := strconv.ParseInt(args[3].Str, 10, 64)
		if err != nil {
			return resp.Value{Typ: resp.ERROR, Str: "ERR invalid expire time"}
		}
		expiresAt = time.Now().UTC().Add(time.Duration(ttl) * time.Second).UnixNano()
	}
	s.Set(key, &MapValue{Value: value, ExpiresAt: expiresAt})
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
	newVal, ok := s.IncrDecr(key, false)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "ERR value is not an integer"}
	}

	return resp.Value{
		Typ: resp.INTEGER,
		Num: newVal,
	}
}

func decrease(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	key := args[0].Str
	newVal, ok := s.IncrDecr(key, true)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "ERR value is not an integer"}
	}

	return resp.Value{
		Typ: resp.INTEGER,
		Num: newVal,
	}
}

func appendStr(s *Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR wrong number of arguments"}
	}

	key := args[0].Str
	val := args[1].Str
	v, ok := s.Get(key)

	if !ok {
		s.Set(key, &MapValue{Value: val, ExpiresAt: 0})
		return resp.Value{Typ: resp.INTEGER, Num: int64(len(val))}
	}

	switch v.Value.(type) {
	case string:
		newStr := v.Value.(string) + val
		n := int64(len(newStr))
		s.Set(key, &MapValue{Value: newStr, ExpiresAt: v.ExpiresAt})
		return resp.Value{Typ: resp.INTEGER, Num: n}
	default:
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}
}

func ping(s *Store, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: resp.STRING, Str: "PONG"}
	}
	return resp.Value{Typ: resp.BULKSTR, Str: args[0].Str}
}
