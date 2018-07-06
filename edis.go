package edis

import (
	"reflect"
	"strings"

	"fmt"

	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/orderedmap"
)

type EventDispatcherInterface interface {
	On(eventName string, callbacks ...interface{})
	OnE(eventName string, callbacks ...interface{}) error
	Trigger(e EventInterface) error
	IsAnyTrigger() bool
	SetAnyTrigger(v bool)
	EnableAnyTrigger()
	DisableAnyTrigger()
	Listeners(key string) (lis *orderedmap.OrderedMap, ok bool)
}

type EventDispatcher struct {
	listeners  map[string]*orderedmap.OrderedMap
	anyTrigger bool
}

func New() *EventDispatcher {
	return &EventDispatcher{}
}

func (ed *EventDispatcher) IsAnyTrigger() bool {
	return ed.anyTrigger
}

func (ed *EventDispatcher) SetAnyTrigger(v bool) {
	ed.anyTrigger = v
}

func (ed *EventDispatcher) EnableAnyTrigger() {
	ed.anyTrigger = true
}

func (ed *EventDispatcher) DisableAnyTrigger() {
	ed.anyTrigger = false
}

func (ed *EventDispatcher) On(eventName string, callbacks ...interface{}) {
	if err := ed.OnE(eventName, callbacks...); err != nil {
		panic(err)
	}
}

func (ed *EventDispatcher) OnE(eventName string, callbacks ...interface{}) error {
	if ed.listeners == nil {
		ed.listeners = make(map[string]*orderedmap.OrderedMap)
	}
	m, ok := ed.listeners[eventName]
	if !ok {
		m = orderedmap.New()
		ed.listeners[eventName] = m
	}
	var (
		ptr uintptr
		ci  CallbackInterface
	)
	for _, cb := range callbacks {
		switch ct := cb.(type) {
		case CallbackInterface:
			ci = ct
		case func(e EventInterface) error:
			ci = CallbackFuncE(ct)
		case func(e EventInterface):
			ci = CallbackFunc(ct)
		default:
			return fmt.Errorf("Invalid Callback type %s", cb)
		}

		switch ciPtr := cb.(type) {
		case interface{ Pointer() uintptr }:
			ptr = ciPtr.Pointer()
		case interface{ PointerOf() interface{} }:
			ptr = reflect.ValueOf(ciPtr.PointerOf()).Pointer()
		default:
			ptr = reflect.ValueOf(ci).Pointer()
		}

		if !m.Has(ptr) {
			m.Set(ptr, ci)
		}
	}
	return nil
}

func (ed *EventDispatcher) Trigger(e EventInterface) (err error) {
	if ed.anyTrigger {
		if err = ed.trigger("*", e); err != nil {
			return
		}

		parts := strings.Split(e.Name(), ":")
		for i, key := range parts[0 : len(parts)-1] {
			key = strings.Join(parts[0:i], ":") + ":" + key + ":*"
			if key[0] == ':' {
				key = key[1:]
			}
			err = ed.trigger(key, e)
			if err != nil {
				return
			}
		}
	}
	return ed.trigger(e.Name(), e)
}

func (ed *EventDispatcher) Listeners(key string) (lis *orderedmap.OrderedMap, ok bool) {
	lis, ok = ed.listeners[key]
	return
}

func (ed *EventDispatcher) trigger(key string, e EventInterface) (err error) {
	oldName := e.CurrentName()
	if oldName != key {
		defer func() {
			e.SetCurrentName(oldName)
		}()
		e.SetCurrentName(key)
	}
	if m, ok := ed.listeners[key]; ok {
		err = m.EachValues(func(value interface{}) error {
			return value.(CallbackInterface).Call(e)
		})
		if err != nil {
			return errwrap.Wrap(err, "Trigger %q", key)
		}
	}
	return
}
