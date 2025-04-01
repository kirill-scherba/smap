// Copyright 2025 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package smap is a thread-safe map with generic any key and any value.
package smap

import "sync"

// Smap is a thread safe map with generic comparable key and any value.
type Smap[K comparable, V any] struct{ sync.Map }

// New returns a pointer to a new Smap.
func New[K comparable, V any]() *Smap[K, V] { return &Smap[K, V]{} }

// Set sets the value for a key.
func (m *Smap[K, V]) Set(key K, value V) {
	m.Map.Store(key, value)
}

// Get returns the value stored in the map for a key, or default V type value if
// no value is present by key.
// The ok result indicates whether value was found in the map.
func (m *Smap[K, V]) Get(key K) (V, bool) {
	if v, ok := m.Map.Load(key); ok {
		if v, ok := v.(V); ok {
			return v, true
		}
	}
	return *new(V), false
}

// Len returns the number of elements in the map.
func (m *Smap[K, V]) Len() (n int) {
	for range m.Map.Range {
		n++
	}
	return
}

// Range calls f sequentially for each key and value present in the map.
func (m *Smap[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
