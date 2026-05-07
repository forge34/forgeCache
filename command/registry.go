// Package command
package command

import (
	"github.com/forge34/forgeCache/resp"
)

type (
	Handler    func([]resp.Value) resp.Value
	HandlerMap map[string]Handler
)

var Handlers = HandlerMap{
	"PING":   ping,
	"SET":    set,
	"DEL":    delete,
	"GET":    get,
	"INCR":   increase,
	"EXISTS": exists,
	"APPEND": append,
	"DECR":   decrease,
}

func delete(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func exists(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func increase(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func get(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func decrease(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func append(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func ping(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "PONG"}
}

func set(args []resp.Value) resp.Value {
	return resp.Value{Typ: resp.STRING, Str: "OK"}
}
