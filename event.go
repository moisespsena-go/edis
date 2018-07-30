package edis

import "fmt"

type EventInterface interface {
	Name() string
	Names() []string
	SetData(value interface{})
	Data() interface{}
	SetResult(value interface{})
	Result() interface{}
	SetError(err error)
	Error() error
	Dispatcher() EventDispatcherInterface
	SetDispatcher(dis EventDispatcherInterface)
	With(name string) (done func())
	WithDispatcher(dis EventDispatcherInterface) (done func())
}

type Event struct {
	names        []string
	PName        string
	PData        interface{}
	PResult      interface{}
	PError       error
	PCurrentName string
	PDispatcher  EventDispatcherInterface
}

func NewEvent(name string, data ...interface{}) *Event {
	e := &Event{PName: name}
	if len(data) == 1 && data[0] != nil {
		e.PData = data[0]
	}
	return e
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{Name=%q, CurrentName=%q, data=%v}", e.PName, e.PCurrentName, e.PData)
}

func (e *Event) Names() []string {
	return e.names
}

func (d *Event) Name() string {
	return d.PName
}

func (d *Event) SetData(value interface{}) {
	d.PData = value
}

func (d *Event) Data() interface{} {
	return d.PData
}

func (d *Event) SetResult(value interface{}) {
	d.PResult = value
}

func (d *Event) Result() interface{} {
	return d.PResult
}

func (d *Event) SetError(err error) {
	d.PError = err
}

func (d *Event) Error() error {
	return d.PError
}

func (d *Event) Dispatcher() EventDispatcherInterface {
	return d.PDispatcher
}

func (d *Event) SetDispatcher(dis EventDispatcherInterface) {
	d.PDispatcher = dis
}

func (e *Event) With(name string) (done func()) {
	oldName := e.PName
	e.names = append(e.names, e.PName)
	e.PName = name
	return func() {
		e.names = e.names[:len(e.names)-1]
		e.PName = oldName
	}
}

func (e *Event) WithDispatcher(dis EventDispatcherInterface) (done func()) {
	old := e.PDispatcher
	e.PDispatcher = dis
	return func() {
		e.PDispatcher = old
	}
}

func EAll(name string) string {
	if name == "" {
		panic("Name is empty")
	}
	return name + ":*"
}
