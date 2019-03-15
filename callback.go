package edis

type Callback interface {
	Call(e EventInterface) error
}

type CallbackFuncE func(e EventInterface) error

func (c CallbackFuncE) Call(e EventInterface) error {
	return c(e)
}

type CallbackFunc func(e EventInterface)

func (c CallbackFunc) Call(e EventInterface) error {
	c(e)
	return nil
}

func SimpleCallback(f func()) Callback {
	return CallbackFunc(func(e EventInterface) {
		f()
	})
}

func SimpleCallbackE(f func() error) Callback {
	return CallbackFuncE(func(e EventInterface) error {
		return f()
	})
}
