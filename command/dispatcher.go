package command

import (
	"strings"

	"github.com/forge34/forgeCache/resp"
)

type Dispatcher struct {
	comamnder *Commander
	store     *Store
}

func NewDispatcher(s *Store) *Dispatcher {
	c := NewCommander()
	return &Dispatcher{comamnder: c, store: s}
}

func (d *Dispatcher) Dispatch(v resp.Value) resp.Value {
	if v.Typ != resp.ARRAY || len(v.Array) == 0 {
		return resp.Value{Typ: resp.ERROR, Str: "ERR invalid command"}
	}

	cmd := strings.ToUpper(v.Array[0].Str)
	args := v.Array[1:]

	handler, ok := d.comamnder.handlers[cmd]
	if !ok {
		return resp.Value{Typ: resp.ERROR, Str: "ERR unknown command '" + cmd + "'"}
	}

	return handler(d.store, args)
}
