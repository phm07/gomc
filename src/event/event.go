package event

import (
	"fmt"
	"reflect"
)

type Bus struct {
	listeners map[reflect.Type][]any
}

func NewBus() *Bus {
	return &Bus{
		listeners: make(map[reflect.Type][]any),
	}
}

func (b *Bus) Emit(data any) error {
	event := reflect.TypeOf(data).Elem()
	for _, listener := range b.listeners[event] {
		ret := reflect.ValueOf(listener).Call([]reflect.Value{reflect.ValueOf(data)})
		err := ret[0]
		if !err.IsNil() {
			return err.Interface().(error)
		}
	}
	return nil
}

func (b *Bus) RegisterListener(listener any) {
	t := reflect.TypeOf(listener)
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("listener must be a function"))
	}
	if t.NumIn() != 1 {
		panic(fmt.Errorf("listener must have exactly one argument"))
	}
	if t.In(0).Kind() != reflect.Ptr {
		panic(fmt.Errorf("listener argument must be a pointer"))
	}
	if t.NumOut() != 1 {
		panic(fmt.Errorf("listener must have exactly one return value"))
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic(fmt.Errorf("listener must return an error"))
	}
	eventType := t.In(0).Elem()
	b.listeners[eventType] = append(b.listeners[eventType], listener)
}
