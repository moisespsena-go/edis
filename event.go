package edis

import "fmt"

type EventInterface interface {
	Name() string
	Data() interface{}
	SetResult(value interface{})
	Result() interface{}
	SetError(err error)
	Error() error
	Parent() EventInterface
	CurrentName() string
	SetCurrentName(value string)
}

type Event struct {
	PName        string
	PData        interface{}
	PResult      interface{}
	PError       error
	PParent      EventInterface
	PCurrentName string
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
	return d.PName
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
