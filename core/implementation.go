package core

import (
	"context"
)

type TransactionType byte

const (
	_ TransactionType = iota
	PUT
	DELETE
)

type Transaction struct {
	Id              int
	TransactionType TransactionType
	Key             string
	Value           string
}

type TransactionLogger interface {
	WritePut(ctx context.Context, key string, value string) error
	WriteDelete(ctx context.Context, key string) error
	ReadAll(ctx context.Context) (chan Transaction, chan error)
}

type Frontend interface {
	Start() error
}

type Encoder interface {
	Encode(string) string
	Decode(string) (string, error)
}
