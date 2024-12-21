package main

import (
	"sync"
)

type KeyValueStore struct {
	Lock sync.RWMutex
	Data map[string]string
}

func NewKvStore() *KeyValueStore {
	return &KeyValueStore{
		sync.RWMutex{}, make(map[string]string)}
}

func (kv *KeyValueStore) SET(k string, v string) {
	kv.Lock.Lock()
	defer kv.Lock.Unlock()
	kv.Data[k] = v
}

func (kv *KeyValueStore) GET(k string) (string, error) {
	kv.Lock.Lock()
	defer kv.Lock.Unlock()
	if output, ok := kv.Data[k]; ok {
		return output, nil
	}
	return "", ErrKeyNotFound
}
