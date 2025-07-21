package core

import (
	"context"
	"errors"
	"sync"
)

var (
	CoreKeyNotFound = errors.New("Could not find any register with this key")
	CorePutError    = errors.New("Error on executing Put")
	CoreDeleteError = errors.New("Error on executing Delete")
)

type KeyValueStore struct {
	Store map[string]string
	m     *sync.RWMutex
	TransactionLogger
	Encoder
}

func NewKeyValueStore(t TransactionLogger, e Encoder) *KeyValueStore {
	kvs := &KeyValueStore{
		TransactionLogger: t,
		Store:             make(map[string]string),
		m:                 &sync.RWMutex{},
		Encoder:           e,
	}

	return kvs
}

func (kvs *KeyValueStore) Put(ctx context.Context, key string, value string) error {
	kvs.m.Lock()
	kvs.Store[key] = value
	kvs.m.Unlock()

	encodedValue := kvs.Encode(value)
	kvs.TransactionLogger.WritePut(ctx, key, encodedValue)

	return nil
}

func (kvs *KeyValueStore) Delete(ctx context.Context, key string) error {
	kvs.m.RLock()
	_, ok := kvs.Store[key]
	kvs.m.RUnlock()

	if !ok {
		return CoreKeyNotFound
	}

	kvs.m.Lock()
	delete(kvs.Store, key)
	kvs.m.Unlock()

	kvs.TransactionLogger.WriteDelete(ctx, key)

	return nil
}

func (kvs KeyValueStore) Get(key string) (string, error) {
	kvs.m.RLock()
	v, ok := kvs.Store[key]
	kvs.m.RUnlock()

	if !ok {
		return "", CoreKeyNotFound
	}

	return v, nil
}

func (kvs *KeyValueStore) Restore(ctx context.Context) error {
	rCh, errCh := kvs.TransactionLogger.ReadAll(ctx)
	var err error

	looper := true
	for looper {
		select {
		case err := <-errCh:
			looper = false
			return err
		case r, ok := <-rCh:
			if !ok {
				looper = false
				break
			}
			switch r.TransactionType {
			case DELETE:
				err := kvs.Delete(ctx, r.Key)
				if err != nil {
					err := errors.Join(err, CorePutError)
					return err
				}
			case PUT:
				decodedValue, err := kvs.Decode(r.Value)
				if err != nil {
					err := errors.Join(err, CorePutError)
					return err
				}

				err = kvs.Put(ctx, r.Key, decodedValue)
				if err != nil {
					err := errors.Join(err, CorePutError)
					return err
				}
			}
		}
	}

	return err
}
