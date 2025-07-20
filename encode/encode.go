package encode

import (
	"errors"
	"github.com/vinmazzi/keyValueStore/core"
)

var (
	EncoderNotValidError = errors.New("Provided encoder is not valid")
)

func NewEncoder(encoderType string) (core.Encoder, error) {
	var encoder core.Encoder

	switch encoderType {
	case "base64":
		encoder = NewBase64Encoder()
	default:
		return nil, EncoderNotValidError
	}

	return encoder, nil
}
