package frontend

import (
	"github.com/vinmazzi/keyValueStore/core"
)

func NewFrontEnd(t string, kvs *core.KeyValueStore) core.Frontend {
	var frontendO core.Frontend
	switch t {
	case "rest":
		frontendO = NewRestFrontend(kvs)
	}

	return frontendO
}

