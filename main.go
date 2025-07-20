package main

import (
	"context"
	"github.com/vinmazzi/keyValueStore/core"
	"github.com/vinmazzi/keyValueStore/frontend"
	"github.com/vinmazzi/keyValueStore/transact"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tl, err := transact.NewTransactionLogger(ctx, os.Getenv("TRANSACTION_LOGGER_BACKEND"))
	if err != nil {
		log.Println(err)
	}

	kvs := core.NewKeyValueStore(tl)
	err = kvs.Restore(ctx)
	if err != nil {
		panic(err)
	}

	fe := frontend.NewFrontEnd(os.Getenv("FRONTEND_TYPE"), kvs)
	err = fe.Start()
	if err != nil {
		panic(err)
	}
}
