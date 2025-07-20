package core

import (
	"context"
	"errors"
	"sync"
)

var KeyNotFound = errors.New("Could not find any register with this key.")

type TransactionLogger interface {
	WritePut(ctx context.Context, key string, value string) error
	WriteDelete(ctx context.Context, key string) error
}

type Frontend interface {
	Start()
}

type KeyValueStore struct {
	Store map[string]string
	m     *sync.RWMutex
	TransactionLogger
}

func (kvs *KeyValueStore) Put(ctx context.Context, key string, value string) error {
	kvs.m.Lock()
	kvs.Store[key] = value
	kvs.m.Unlock()

	kvs.TransactionLogger.WritePut(ctx, key, value)

	return nil
}

func (kvs *KeyValueStore) Delete(ctx context.Context, key string) error {
	kvs.m.RLock()
	_, ok := kvs.Store[key]
	kvs.m.RUnlock()

	if !ok {
		return KeyNotFound
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
		return "", KeyNotFound
	}

	return v, nil
}
