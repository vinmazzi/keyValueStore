package frontend

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vinmazzi/keyValueStore/core"
)

var (
	RestHttpServerError         = errors.New("Http server failed on start")
	RestHttpIntConversionError  = errors.New("Error on the timeout conversion")
	RestHttpErrorReadingReqBody = errors.New("Error reading RequestBody")
	RestHttpErrorPutFunction    = errors.New("Error When calling the Put Method from core")
	RestHttpErrorDeleteFunction = errors.New("Error When calling the Delete Method from core")
)

type RestFrontEnd struct {
	*core.KeyValueStore
}

func NewRestFrontend(kvs *core.KeyValueStore) *RestFrontEnd {
	restFrontEnd := &RestFrontEnd{
		KeyValueStore: kvs,
	}

	log.Println("Starting REST frontend")
	return restFrontEnd
}

func (rfe RestFrontEnd) GetKeyHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	v, ok := rfe.Store[key]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(v))
}

func (rfe RestFrontEnd) PutKeyHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		err = errors.Join(err, RestHttpErrorReadingReqBody)
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = rfe.Put(ctx, key, string(body))
	if err != nil {
		err = errors.Join(err)
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusCreated)
}

func (rfe RestFrontEnd) DeleteKeyHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := rfe.Delete(ctx, key)
	if err != nil {
		err = errors.Join(err, RestHttpErrorDeleteFunction)
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusCreated)
}

func (rfe *RestFrontEnd) Start() error {

	router := mux.NewRouter()
	router.HandleFunc("/v1/key/{key}", rfe.GetKeyHandler).Methods("GET")
	router.HandleFunc("/v1/key/{key}", rfe.PutKeyHandler).Methods("PUT")
	router.HandleFunc("/v1/key/{key}", rfe.DeleteKeyHandler).Methods("DELETE")

	secondsTimeout, err := strconv.Atoi(os.Getenv("HTTP_IDLE_TIMEOUT"))
	if err != nil {
		return errors.Join(err, RestHttpIntConversionError)
	}

	server := http.Server{
		Addr:        os.Getenv("HTTP_LISTEN_PORT"),
		IdleTimeout: time.Duration(secondsTimeout) * time.Second,
		Handler:     router,
	}

	err = server.ListenAndServe()
	if err != nil {
		err = errors.Join(err, RestHttpServerError)
		panic(err)
	}

	return err
}
