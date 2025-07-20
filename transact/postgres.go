package transact

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/vinmazzi/keyValueStore/core"
	"log"
)

type PostgresConnParam struct {
	host     string
	user     string
	password string
	sslmode  string
	dbname   string
}

type PostgresTransactionLogger struct {
	*sql.DB
}

var (
	PostgresErrorOpeningConnection = errors.New("Error Opening connection with postgres.")
	PostgresConnectionError        = errors.New("Error connecting with Postgres")
	//TODO: This needs to be fixed at some point
	PostgresConnNotInsecure           = errors.New("We only accept insecure connections for now")
	PostgresErrorExecutingPutQuery    = errors.New("Was not able to execute the Put WritePut query")
	PostgresErrorExecutingDeleteQuery = errors.New("Could not execute the Delete query on WriteDelete")
	PostgresErrorExecutingSectQuery   = errors.New("Could not execute the Select Query on ReadAll")
)

const (
	transactionsTable = "transactions"
)

func (cp PostgresConnParam) CheckSslDisabled() bool {
	if cp.sslmode != "disable" {
		return false
	}
	return true
}

func NewPostgresTransactionLogger(ctx context.Context, connParams PostgresConnParam) (*PostgresTransactionLogger, error) {
	var ptl PostgresTransactionLogger
	var err error

	ok := connParams.CheckSslDisabled()
	if !ok {
		return nil, PostgresConnNotInsecure
	}

	connString := fmt.Sprintf("host=%s user=%s password=%s sslmode=%s dbname=%s", connParams.host, connParams.user, connParams.password, connParams.sslmode, connParams.dbname)
	ptl.DB, err = sql.Open("postgres", connString)
	if err != nil {
		return nil, errors.Join(err, PostgresErrorOpeningConnection)
	}

	err = ptl.Ping()
	if err != nil {
		return nil, errors.Join(err, PostgresConnectionError)
	}

	log.Println("Starting postgresl transaction logger.")

	return &ptl, err
}

func TableCheck() error { return nil }

func (ptl *PostgresTransactionLogger) WritePut(ctx context.Context, key string, value string) error {
	query := `INSERT INTO transactions (type, key, value) VALUES ($1, $2, $3)`

	_, err := ptl.DB.ExecContext(ctx, query, core.PUT, &key, &value)
	if err != nil {
		errors.Join(err, PostgresErrorExecutingPutQuery)
		return err
	}

	return nil
}

func (ptl *PostgresTransactionLogger) WriteDelete(ctx context.Context, key string) error {
	query := `INSERT INTO transactions (type, key) VALUES ($1, $2)`

	_, err := ptl.DB.ExecContext(ctx, query, core.DELETE, &key)
	if err != nil {
		return errors.Join(err, PostgresErrorExecutingDeleteQuery)
	}

	return nil
}

func (ptl *PostgresTransactionLogger) ReadAll(ctx context.Context) (chan core.Transaction, chan error) {
	errCh := make(chan error)
	trCh := make(chan core.Transaction)

	go func() {
		defer close(errCh)
		defer close(trCh)

		query := `SELECT * FROM transactions`
		result, err := ptl.DB.QueryContext(ctx, query)
		if err != nil {
			log.Println("There is an error here:", err)
			errCh <- errors.Join(err, PostgresErrorExecutingSectQuery)
		}
		defer result.Close()

		for result.Next() {
			t := core.Transaction{}
			err = result.Scan(&t.Id, &t.TransactionType, &t.Key, &t.Value)
			if err != nil {
				log.Println("The error is: ", err)
				return
			}
			trCh <- t
		}
	}()

	return trCh, errCh
}
