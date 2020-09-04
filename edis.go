package edis

import (
	"strings"

	"fmt"

	"github.com/moisespsena-go/logging"
)

type EventDispatcherInterface interface {
	On(eventName string, callbacks ...interface{})
	OnE(eventName string, callbacks ...interface{}) error
	Trigger(e EventInterface) error
	IsAnyTrigger() bool
	SetAnyTrigger(v bool)
	EnableAnyTrigger()
	DisableAnyTrigger()
	Listeners(key string) (lis []Callback, ok bool)
	AllListeners() map[string][]Callback
	SetLogger(log logging.Logger)
	Logger() logging.Logger
	SetDebug(v bool)
	IsDebugEnabled() bool
	EnableDebug()
	DisableDebug()
	Dispatcher() EventDispatcherInterface
	SetDispatcher(dis EventDispatcherInterface)
}

type EventDispatcher struct {
	dispatcher EventDispatcherInterface
	listeners  map[string][]Callback
	anyTrigger bool
	log        logging.Logger
	debug      bool
	debugFunc  func(dis EventDispatcherInterface, key string, e EventInterface)
}

func New() *EventDispatcher {
	return &EventDispatcher{}
}

func (ed *EventDispatcher) SetLogger(log logging.Logger) {
	ed.log = log
}

func (ed *EventDispatcher) Logger() logging.Logger {
	return ed.log
}
func (ed *EventDispatcher) IsDebugEnabled() bool {
	return ed.debug
}

func (ed *EventDispatcher) SetDebug(v bool) {
	ed.debug = v
	if v {
		if ed.log == nil {
			ed.log = log
		}
		if ed.debugFunc == nil {
			ed.debugFunc = DefaultEventDebug
		}
	}
}

func (ed *EventDispatcher) EnableDebug() {
	ed.SetDebug(true)
}

func (ed *EventDispatcher) DisableDebug() {
	ed.debug = false
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
		ed.listeners = make(map[string][]Callback)
	}
	m, ok := ed.listeners[eventName]
	if !ok {
		m = make([]Callback, 0, len(callbacks))
	}
	var ci Callback
	for _, cb := range callbacks {
		switch ct := cb.(type) {
		case Callback:
			ci = ct
		case func(e EventInterface) error:
			ci = CallbackFuncE(ct)
		case func(e EventInterface):
			ci = CallbackFunc(ct)
		case func():
			ci = SimpleCallback(ct)
		case func() error:
			ci = SimpleCallbackE(ct)
		default:
			return fmt.Errorf("Invalid Callback type %s", cb)
		}

		m = append(m, ci)
	}
	ed.listeners[eventName] = m
	return nil
}

func (ed *EventDispatcher) GetDefinedDispatcher() EventDispatcherInterface {
	return ed.dispatcher
}

func (ed *EventDispatcher) Dispatcher() EventDispatcherInterface {
	if ed.dispatcher == nil {
		return ed
	}

	return ed.dispatcher
}

func (ed *EventDispatcher) SetDispatcher(dis EventDispatcherInterface) {
	ed.dispatcher = dis
}

func (ed *EventDispatcher) Trigger(e EventInterface) (err error) {
	defer e.WithDispatcher(ed.Dispatcher())()
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

func (ed *EventDispatcher) Listeners(key string) (lis []Callback, ok bool) {
	lis, ok = ed.listeners[key]
	return
}

func (ed *EventDispatcher) AllListeners() map[string][]Callback {
	return ed.listeners
}

func (ed *EventDispatcher) trigger(key string, e EventInterface) (err error) {
	if ed.debug {
		ed.debugFunc(ed, key, e)
	}
	defer e.With(key)()
	if m, ok := ed.listeners[key]; ok {
		for _, cb := range m {
			if err = cb.Call(e); err != nil {
				return
			}
		}
	}
	return
}
