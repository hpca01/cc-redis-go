package main

import (
	"sync"
	"time"
)

var maxTime = time.Unix(1<<63-62135596801, 999999999)

type Value struct {
	value string
	timer time.Time
}

func (v *Value) isExpired() bool {
	return v.timer.After(time.Now())
}

func NewValue(v string, time *time.Time) Value {
	if time == nil {
		//max time
		return Value{v, maxTime}
	}
	return Value{v, *time}
}

type KeyValueStore struct {
	lock sync.RWMutex
	data map[string]Value
}

func NewKvStore() *KeyValueStore {
	return &KeyValueStore{
		sync.RWMutex{}, make(map[string]Value)}
}

func (kv *KeyValueStore) SET(k string, v Value) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[k] = v
}

func (kv *KeyValueStore) GET(k string) (string, error) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	if output, ok := kv.data[k]; ok {
		if output.isExpired() {
			delete(kv.data, k)
			return "", ErrKeyExpired
		}
		return output.value, nil
	}
	return "", ErrKeyNotFound
}
