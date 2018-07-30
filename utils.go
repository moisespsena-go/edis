package edis

import (
	"fmt"
)

func DefaultEventDebug(dis EventDispatcherInterface, key string, e EventInterface) {
	dis.Logger().Debug(fmt.Sprintf("Trigger %q -> %s", key, e))
}
