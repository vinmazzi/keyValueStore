package encode

import (
	"encoding/base64"
	"errors"
)

var (
	Base64DecodeError = errors.New("Could not decode srting on method Decode")
)

type Base64Encoder struct {
	base64.Encoding
}

func NewBase64Encoder() *Base64Encoder {
	b64 := &Base64Encoder{
		Encoding: *base64.StdEncoding,
	}

	return b64
}

func (b64 *Base64Encoder) Encode(s string) string {
	return b64.EncodeToString([]byte(s))
}

func (b64 *Base64Encoder) Decode(s string) (string, error) {
	decoded, err := b64.DecodeString(s)
	if err != nil {
		err = errors.Join(err, Base64DecodeError)
		return "", err
	}

	return string(decoded), nil
}
