# go-edis
Go Lang Simples Event Dispatcher API and Implemetation

## Installation

```bash
go get -u github.com/moisespsena/go-edis
```

## Usage

```go
package main

import (
	"fmt"
	"strings"

	"github.com/moisespsena/go-edis"
)

func main() {
	pad := func(v string) string {
		return v + strings.Repeat(" ", 11-len(v))
	}
	ed := edis.New()
	ed.On("*", func(e edis.EventInterface) error {
		fmt.Println(pad("0 -> *"), e)
		return nil
	})
	ed.On("event", func(e edis.EventInterface) error {
		fmt.Println(pad("1 -> event"), e)
		return nil
	})
	ed.On("a:*", func(e edis.EventInterface) error {
		fmt.Println(pad("2 -> a:*"), e)
		return nil
	})
	ed.On("a:b:*", func(e edis.EventInterface) error {
		fmt.Println(pad("3 -> a:b:*"), e)
		return nil
	})
	ed.On("a:b:c", func(e edis.EventInterface) error {
		fmt.Println(pad("4 -> a:b:c"), e)
		return nil
	})

	trigger := func() {
		ed.Trigger(edis.NewEvent("event", "data"))
		ed.Trigger(edis.NewEvent("a:b:c", "abc"))
	}

	fmt.Println("----- any trigger disabled -----")
	trigger()
	
	fmt.Println("----- any trigger enabled -----")
	ed.EnableAnyTrigger()
	trigger()
}
```

Result:

```
----- any trigger disabled -----
1 -> event  Event{Name="event", CurrentName="event", data=data}
4 -> a:b:c  Event{Name="a:b:c", CurrentName="a:b:c", data=abc}
----- any trigger enabled -----
0 -> *      Event{Name="event", CurrentName="*", data=data}
1 -> event  Event{Name="event", CurrentName="event", data=data}
0 -> *      Event{Name="a:b:c", CurrentName="*", data=abc}
2 -> a:*    Event{Name="a:b:c", CurrentName="a:*", data=abc}
3 -> a:b:*  Event{Name="a:b:c", CurrentName="a:b:*", data=abc}
4 -> a:b:c  Event{Name="a:b:c", CurrentName="a:b:c", data=abc}
```