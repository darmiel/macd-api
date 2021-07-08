package common

import "sync"

type AtomicInt64 struct {
	value int64
	mu    sync.Mutex
}

func (a *AtomicInt64) Add(n int64) int64 {
	a.mu.Lock()
	a.value += n
	val := a.value
	a.mu.Unlock()
	return val
}
func (a *AtomicInt64) Get() int64 {
	a.mu.Lock()
	val := a.value
	a.mu.Unlock()
	return val
}
