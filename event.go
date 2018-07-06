package edis

import "fmt"

type EventInterface interface {
	Name() string
	SetData(value interface{})
	Data() interface{}
	SetResult(value interface{})
	Result() interface{}
	SetError(err error)
	Error() error
	SetParent(parent EventInterface)
	Parent() EventInterface
	CurrentName() string
	SetCurrentName(value string)
	Dispatcher() EventDispatcherInterface
	SetDispatcher(dis EventDispatcherInterface)
}

type Event struct {
	PName        string
	PData        interface{}
	PResult      interface{}
	PError       error
	PParent      EventInterface
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

func (d *Event) CurrentName() string {
	if d.PCurrentName == "" {
		return d.PName
	}
	return d.PCurrentName
}

func (d *Event) SetCurrentName(value string) {
	d.PCurrentName = value
}

func (d *Event) Name() string {
	if d.PName == "" && d.PParent != nil {
		return d.PParent.Name()
	}
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
	if d.PParent != nil {
		d.PParent.SetResult(value)
	}
}

func (d *Event) Result() interface{} {
	return d.PResult
}

func (d *Event) SetError(err error) {
	d.PError = err
	if d.PParent != nil {
		d.PParent.SetError(err)
	}
}

func (d *Event) Error() error {
	return d.PError
}

func (d *Event) Parent() EventInterface {
	return d.PParent
}

func (d *Event) SetParent(parent EventInterface) {
	d.PParent = parent
}

func (d *Event) Dispatcher() EventDispatcherInterface {
	return d.PDispatcher
}

func (d *Event) SetDispatcher(dis EventDispatcherInterface) {
	d.PDispatcher = dis
}

func EAll(name string) string {
	if name == "" {
		panic("Name is empty")
	}
	return name + ":*"
}
