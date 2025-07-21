package transact

import (
	"context"
	"errors"
	"github.com/vinmazzi/keyValueStore/core"
	"os"
)

var (
	TransactionLoggerNotSupportedError           = errors.New("The requested transaction logger is not supported.")
	TransactionLoggerPostgresInitializationError = errors.New("Error on creating a Postgres transactionLogger")
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
			err = errors.Join(err, TransactionLoggerPostgresInitializationError)
			return nil, err
		}
	default:
		return nil, TransactionLoggerNotSupportedError
	}

	return logger, nil
}
