package edis

type CallbackInterface interface {
	Call(e EventInterface) error
}

type CallbackFunc func(e EventInterface) error

func (c CallbackFunc) Call(e EventInterface) error {
	return c(e)
}
