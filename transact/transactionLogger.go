package transact

import (
	"context"
	"github.com/vinmazzi/keyValueStore/core"
	"os"
)

type TransactionType byte

type Transaction struct {
	id              int
	transactionType TransactionType
	key             string
	value           string
}

const (
	_ TransactionType = iota
	PUT
	DELETE
)

func NewTransactionLogger(ctx context.Context, loggerType string) (core.TransactionLogger, error) {
	var logger core.TransactionLogger
	var err error
	switch loggerType {
	case "postgres":
		connParams := PostgresConnParam{
			host:     os.Getenv("POSTGRES_HOST"),
			user:     os.Getenv("POSTGRES_USERNAME"),
			password: os.Getenv("POSTGRES_PASSWORD"),
			sslmode:  os.Getenv("POSTGRES_SSLMODE"),
			dbname:   os.Getenv("POSTGRES_DATABASE"),
		}

		logger, err = NewPostgresTransactionLogger(ctx, connParams)
		if err != nil {
			return nil, err
		}
	}

	return logger, nil
}
