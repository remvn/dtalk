package cmap

import (
	"iter"
	"sync"
)

type CMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func New[K comparable, V any]() *CMap[K, V] {
	cmap := &CMap[K, V]{
		m:     make(map[K]V),
		mutex: sync.RWMutex{},
	}
	return cmap
}

func (cmap *CMap[K, V]) Get(key K) V {
	cmap.mutex.RLock()
	defer cmap.mutex.RUnlock()
	return cmap.m[key]
}

func (cmap *CMap[K, V]) Load(key K) (V, bool) {
	cmap.mutex.RLock()
	defer cmap.mutex.RUnlock()
	v, ok := cmap.m[key]
	return v, ok
}

func (cmap *CMap[K, V]) Set(key K, val V) {
	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	cmap.m[key] = val
}

func (cmap *CMap[K, V]) Delete(key K) {
	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	delete(cmap.m, key)
}

// invocation of other operations like Set, Get, Delete...
// inside yield (for block body) will create a deadlock
func (cmap *CMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(key K, val V) bool) {
		cmap.mutex.Lock()
		defer cmap.mutex.Unlock()
		for k, v := range cmap.m {
			if !yield(k, v) {
				break
			}
		}
	}
}
