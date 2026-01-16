package pola

import (
	"errors"
	"fmt"
	"maps"
	"sync"
)

var (
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrEntryDoesNotExists = errors.New("entry does not exists")
)

// Registry for storing entry.
type Registry[K comparable, V any] interface {
	Register(k K, v V) error
	MustRegister(k K, v V)
	Exists(k K) bool
	Set(k K, v V)
	Get(k K) (V, error)
	MustGet(k K) V
	Map() map[K]V
}

type mapRegistry[K comparable, V any] map[K]V

// NewRegistry create map-based registry.
// Please note that this registry is not safe for concurrent usage.
func NewRegistry[K comparable, V any]() Registry[K, V] {
	return make(mapRegistry[K, V])
}

func (m mapRegistry[K, V]) Map() map[K]V {
	if len(m) == 0 {
		return nil
	}
	d := make(map[K]V)
	maps.Copy(d, m)

	return d
}

func (m mapRegistry[K, V]) Set(k K, v V) {
	m[k] = v
}

func (m mapRegistry[K, V]) Register(k K, v V) error {
	if _, ok := m[k]; ok {
		return ErrDuplicateEntry
	}
	m[k] = v
	return nil
}

func (m mapRegistry[K, V]) MustRegister(k K, v V) {
	if _, ok := m[k]; ok {
		panic(fmt.Sprintf("duplicate entry `%v`", k))
	}
	m[k] = v
}

func (m mapRegistry[K, V]) Exists(k K) bool {
	_, ok := m[k]
	return ok
}
func (m mapRegistry[K, V]) Get(k K) (V, error) {
	v, ok := m[k]
	if ok {
		return v, nil
	}
	return v, fmt.Errorf("key: %v, %w", k, ErrEntryDoesNotExists)
}
func (m mapRegistry[K, V]) MustGet(k K) V {
	v, ok := m[k]
	if !ok {
		panic(fmt.Sprintf("entry `%v` does not exists", k))
	}
	return v
}

type syncMapRegistry[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

// NewRegistry create map-based registry guarded with Mutex.
// This registry is safe for concurrent usage.
func NewSyncRegistry[K comparable, V any]() Registry[K, V] {
	r := &syncMapRegistry[K, V]{
		m: make(map[K]V),
	}
	return r
}

func (r *syncMapRegistry[K, V]) Map() map[K]V {
	r.RLock()
	defer r.RUnlock()

	if len(r.m) == 0 {
		return nil
	}
	d := make(map[K]V)
	maps.Copy(d, r.m)

	return d
}

func (r *syncMapRegistry[K, V]) Set(k K, v V) {
	r.Lock()
	defer r.Unlock()

	r.m[k] = v
}

func (r *syncMapRegistry[K, V]) Register(k K, v V) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[k]; ok {
		return ErrDuplicateEntry
	}
	r.m[k] = v
	return nil
}

func (r *syncMapRegistry[K, V]) MustRegister(k K, v V) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[k]; ok {
		panic(fmt.Sprintf("duplicate entry `%v`", k))
	}
	r.m[k] = v
}

func (r *syncMapRegistry[K, V]) Exists(k K) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.m[k]
	return ok
}
func (r *syncMapRegistry[K, V]) Get(k K) (V, error) {
	r.RLock()
	defer r.RUnlock()

	v, ok := r.m[k]
	if ok {
		return v, nil
	}
	return v, fmt.Errorf("key: %v, %w", k, ErrEntryDoesNotExists)
}
func (r *syncMapRegistry[K, V]) MustGet(k K) V {
	r.RLock()
	defer r.RUnlock()

	v, ok := r.m[k]
	if !ok {
		panic(fmt.Sprintf("entry `%v` does not exists", k))
	}
	return v
}
