// Speed test smap and original map.

package smap

import (
	"fmt"
	"sync"
	"testing"
)

// Build-in map with mutex
type Bmap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// New returns a pointer to a new Bmap.
func NewBmap[K comparable, V any]() *Bmap[K, V] {
	return &Bmap[K, V]{m: make(map[K]V)}
}

// Set sets the value for a key.
func (m *Bmap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}

// Get returns the value stored in the map for a key, or default V type value if
// no value is present by key.
// The ok result indicates whether value was found in the map.
func (m *Bmap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

// MapInterface is a map interface for benchmark test.
type MapInterface[K, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
}

// Set sets the value for a key for benchmark test.
func set[K, V any](m MapInterface[K, V], wg *sync.WaitGroup, key K, value V) {
	m.Set(key, value)
	wg.Done()
}

// Get returns the value stored in the map for a key, or default V type value if
// no value is present by key for benchmark test.
// The ok result indicates whether value was found in the map.
func get[K, V any](m MapInterface[K, V], wg *sync.WaitGroup, key K) (V, bool) {
	v, ok := m.Get(key)
	wg.Done()
	return v, ok
}

// BenchmarkSmap benchmark test for smap.
func BenchmarkSmap(b *testing.B) {
	m := New[string, int]()
	wg := sync.WaitGroup{}
	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("test%d", i)
		wg.Add(2)
		go set(m, &wg, key, i)
		go get(m, &wg, key)
	}
	wg.Wait()
}

// BenchmarkBMap benchmark test for build-in map.
func BenchmarkBMap(b *testing.B) {
	m := NewBmap[string, int]()
	wg := sync.WaitGroup{}
	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("test%d", i)
		wg.Add(2)
		go set(m, &wg, key, i)
		go get(m, &wg, key)
	}
	wg.Wait()
}

// TestCompareBenchmarks compares the performance of smap and build-in map.
func TestCompareBenchmarks(t *testing.T) {

	// Run benchmarks
	b1 := testing.Benchmark(BenchmarkSmap)
	b2 := testing.Benchmark(BenchmarkBMap)

	// Print results
	t.Log(b1)
	t.Log(b2)

	// Check results ns per operation are correct
	if b1.NsPerOp() > b2.NsPerOp() {
		t.Errorf("Smap is slower than BMap: %d ns/op vs %d ns/op",
			b1.NsPerOp(), b2.NsPerOp())
	}

	// Check results alloced bytes per operation are correct
	if b1.AllocedBytesPerOp() > b2.AllocedBytesPerOp() {
		t.Errorf("Smap allocates more bytes per operation than BMap: %d B/op vs %d B/op",
			b1.AllocedBytesPerOp(), b2.AllocedBytesPerOp())
	}

	// Print comparison results
	t.Logf(
		"Smap is %.2f times faster than BMap",
		float64(b1.NsPerOp())/float64(b2.NsPerOp()),
	)
	t.Logf(
		"Smap allocates %.2f%% less bytes per operation than BMap",
		100-100*float64(b1.AllocedBytesPerOp())/float64(b2.AllocedBytesPerOp()),
	)
}
