package metrics

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// NewFloatCounter registers and returns new counter of float64 type with the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned counter is safe to use from concurrent goroutines.
func NewFloatCounter(name string) *FloatCounter {
	return defaultSet.NewFloatCounter(name)
}

// FloatCounter is a float64 counter guarded by RWmutex.
//
// It may be used as a gauge if Add and Sub are called.
type FloatCounter struct {
	mu            sync.Mutex
	n             float64
	lastWriteTime int64
}

// Add adds n to fc.
func (fc *FloatCounter) Add(n float64) {
	fc.mu.Lock()
	fc.n += n
	fc.lastWriteTime = time.Now().Unix()
	fc.mu.Unlock()
}

// Sub substracts n from fc.
func (fc *FloatCounter) Sub(n float64) {
	fc.mu.Lock()
	fc.n -= n
	fc.lastWriteTime = time.Now().Unix()
	fc.mu.Unlock()
}

// Get returns the current value for fc.
func (fc *FloatCounter) Get() float64 {
	fc.mu.Lock()
	n := fc.n
	fc.lastWriteTime = time.Now().Unix()
	fc.mu.Unlock()
	return n
}

// Set sets fc value to n.
func (fc *FloatCounter) Set(n float64) {
	fc.mu.Lock()
	fc.n = n
	fc.lastWriteTime = time.Now().Unix()
	fc.mu.Unlock()
}

// marshalTo marshals fc with the given prefix to w.
func (fc *FloatCounter) marshalTo(prefix string, w io.Writer) {
	v := fc.Get()
	fmt.Fprintf(w, "%s %g\n", prefix, v)
}

func (c *FloatCounter) getLastWriteTime() int64 {
	return atomic.LoadInt64(&c.lastWriteTime)
}

func (fc *FloatCounter) metricType() string {
	return "counter"
}

// GetOrCreateFloatCounter returns registered FloatCounter with the given name
// or creates new FloatCounter if the registry doesn't contain FloatCounter with
// the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned FloatCounter is safe to use from concurrent goroutines.
//
// Performance tip: prefer NewFloatCounter instead of GetOrCreateFloatCounter.
func GetOrCreateFloatCounter(name string) *FloatCounter {
	return defaultSet.GetOrCreateFloatCounter(name)
}
