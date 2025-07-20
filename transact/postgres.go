package transact

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
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
	PostgresErrorExecutingDeleteQuery = errors.New("Could not execute the Delete query on WriteDelete.")
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

	return &ptl, err
}

func TableCheck() error { return nil }

func (ptl *PostgresTransactionLogger) WritePut(ctx context.Context, key string, value string) error {

	query := "INSERT INTO $1 (type, key, value) VALUES ($2, $3, $4)"

	_, err := ptl.DB.ExecContext(ctx, query, transactionsTable, PUT, &key, &value)
	if err != nil {
		errors.Join(err, PostgresErrorExecutingPutQuery)
		return err
	}

	return nil
}
func (ptl *PostgresTransactionLogger) WriteDelete(ctx context.Context, key string) error {
	query := "INSERT INTO $1 (type, key) VALUES ($2, $3)"

	_, err := ptl.DB.ExecContext(ctx, query, transactionsTable, &key)
	if err != nil {
		return errors.Join(err, PostgresErrorExecutingDeleteQuery)
	}

	return nil
}
