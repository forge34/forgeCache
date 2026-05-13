package command

import (
	"github.com/forge34/forgeCache/resp"
)

type Dispatcher struct {
	comamnder *Commander
}

func NewDispatcher(s *Store) *Dispatcher {
	c := NewCommander(s)
	return &Dispatcher{comamnder: c}
}

func (d *Dispatcher) Dispatch(v resp.Value) resp.Value {
	result := d.comamnder.Execute(v)

	return result
}
