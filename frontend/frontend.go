package frontend

import (
	"github.com/vinmazzi/keyValueStore/core"
)

func NewFrontEnd(t string, kvs *core.KeyValueStore) core.Frontend {
	var frontendObj core.Frontend
	switch t {
	case "rest":
		frontendObj = NewRestFrontend(kvs)
	}

	return frontendObj
}

