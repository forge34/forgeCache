// Package command
package command

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/forge34/forgeCache/resp"
)

type Handler struct {
	fn       func(*Store, []resp.Value) resp.Value
	writeCmd bool
}

type (
	HandlerMap map[string]Handler
)

const AOFPath = "appendonly.aof"

type Commander struct {
	handlers HandlerMap
	store    *Store
	aof      *os.File
}

func NewCommander(s *Store) *Commander {
	file, err := os.OpenFile("appendonly.aof", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	return &Commander{
		aof:   file,
		store: s,
		handlers: HandlerMap{
			"PING":    {fn: ping, writeCmd: false},
			"SET":     {fn: set, writeCmd: true},
			"DEL":     {fn: deleteKey, writeCmd: true},
			"GET":     {fn: get, writeCmd: false},
			"INCR":    {fn: increase, writeCmd: true},
			"DECR":    {fn: decrease, writeCmd: true},
			"EXISTS":  {fn: exists, writeCmd: false},
			"APPEND":  {fn: appendStr, writeCmd: true},
			"EXPIRE":  {fn: expire, writeCmd: true},
			"TTL":     {fn: checkTTL, writeCmd: false},
			"HSET":    {fn: hset, writeCmd: true},
			"HGET":    {fn: hget, writeCmd: false},
			"HDEL":    {fn: hdel, writeCmd: true},
			"HGETALL": {fn: hgetall, writeCmd: false},
			"HEXISTS": {fn: hexists, writeCmd: false},
			"HLEN":    {fn: hlen, writeCmd: false},
		},
	}
}

func (c *Commander) Execute(v resp.Value , write bool) resp.Value {
	if v.Typ != resp.ARRAY || len(v.Array) == 0 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR invalid command"}
	}

	cmd := strings.ToUpper(v.Array[0].Str)
	args := v.Array[1:]

	commandStr := cmd

	for _, v := range args {
		commandStr += " " + v.Str
	}

	handler, ok := c.handlers[cmd]

	if handler.writeCmd && write {

		_, err := c.aof.Write(v.Marshal())
		if err != nil {
			log.Println("AOF write error:", err)
		}
	}
	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "ERR unknown command '" + cmd + "'"}
	}

	return handler.fn(c.store, args)
}

func checkTTL(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{
			Typ: resp.ERROR,
			Str: "ERR wrong number of arguments",
		}
	}

	key := args[0].Str

	entry, ok := s.Get(key)
	if !ok {
		return resp.Value{
			Typ: resp.INTEGER,
			Num: -2,
		}
	}

	if entry.ExpiresAt == 0 {
		return resp.Value{
			Typ: resp.INTEGER,
			Num: -1,
		}
	}

	timeLeft := entry.ExpiresAt - time.Now().UnixNano()

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

func hlen(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str

	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.INTEGER, Num: 0}
	}

	hsh, ok := entry.Value.(map[string]string)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	return resp.Value{Typ: resp.INTEGER, Num: int64(len(hsh))}
}

func hexists(s *Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str

	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.INTEGER, Num: 0}
	}

	hsh, ok := entry.Value.(map[string]string)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	_, ok = hsh[args[1].Str]

	if !ok {
		return resp.Value{Typ: resp.INTEGER, Num: 0}
	}

	return resp.Value{Typ: resp.INTEGER, Num: 1}
}

func hgetall(s *Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str

	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.ARRAY, Array: []resp.Value{}}
	}

	hsh, ok := entry.Value.(map[string]string)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	var arr []resp.Value

	for field, val := range hsh {
		arr = append(arr, resp.Value{Typ: resp.BULKSTR, Str: field})
		arr = append(arr, resp.Value{Typ: resp.BULKSTR, Str: val})
	}

	return resp.Value{Typ: resp.ARRAY, Array: arr}
}

func hdel(s *Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str

	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.BULKSTR, IsNil: true}
	}

	hsh, ok := entry.Value.(map[string]string)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	var deleted int64 = 0
	for _, field := range args[1:] {
		_, ok := hsh[field.Str]
		if ok {
			delete(hsh, field.Str)
			deleted++
		}
	}

	if len(hsh) == 0 {
		s.Delete(key)
	}

	return resp.Value{Typ: resp.INTEGER, Num: deleted}
}

func hget(s *Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str
	entry, ok := s.Get(key)

	if !ok {
		return resp.Value{Typ: resp.BULKSTR, IsNil: true}
	}

	hsh, ok := entry.Value.(map[string]string)

	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	value := args[1].Str
	v, ok := hsh[value]
	if !ok {
		return resp.Value{Typ: resp.BULKSTR, IsNil: true}
	}

	return resp.Value{Typ: resp.BULKSTR, Str: v}
}

func hset(s *Store, args []resp.Value) resp.Value {
	if len(args) < 3 {
		return resp.Value{Typ: resp.ERROR, Str: "Too few arguments"}
	}

	key := args[0].Str
	entry, ok := s.Get(key)

	var hash map[string]string

	if !ok {
		hash = make(map[string]string)
		s.Set(key, &MapValue{Value: hash})
	} else {
		hash, ok = entry.Value.(map[string]string)
		if !ok {
			return resp.Value{Typ: resp.ERROR, Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
	}
	if len(args[1:])%2 != 0 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR syntax error"}
	}
	var added int64
	for i := 1; i+1 < len(args); i += 2 {
		field := args[i].Str
		value := args[i+1].Str
		if _, exists := hash[field]; !exists {
			added++
		}
		hash[field] = value
	}

	return resp.Value{Typ: resp.INTEGER, Num: added}
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

	keys := make([]string, 0, len(args))

	for _, arg := range args {
		keys = append(keys, arg.Str)
	}

	count := s.Exists(keys)
	return resp.Value{Typ: resp.INTEGER, Num: count}
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
	entry, ok := s.Get(key)

	if !ok {
		s.Set(key, &MapValue{Value: val, ExpiresAt: 0})
		return resp.Value{Typ: resp.INTEGER, Num: int64(len(val))}
	}

	switch entry.Value.(type) {
	case string:
		newStr := entry.Value.(string) + val
		n := int64(len(newStr))
		s.Set(key, &MapValue{Value: newStr, ExpiresAt: entry.ExpiresAt})
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
