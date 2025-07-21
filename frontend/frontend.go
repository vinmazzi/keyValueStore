package frontend

import (
	"errors"
	"github.com/vinmazzi/keyValueStore/core"
)

var (
	FrontendUnavailableError = errors.New("The required frontend does not exist")
)

func NewFrontEnd(t string, kvs *core.KeyValueStore) (core.Frontend, error) {
	var frontendObj core.Frontend
	switch t {
	case "rest":
		frontendObj = NewRestFrontend(kvs)
	default:
		return nil, FrontendUnavailableError
	}

	return frontendObj, nil
}
