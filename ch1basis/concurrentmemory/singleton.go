package concurrentmemory

import (
	"sync"
	"sync/atomic"
)

type singleton struct{}

var (
	instance    *singleton
	initialized uint32
	mutex       sync.Mutex
	once        sync.Once
)

// 原子操作配合互斥锁可以实现非常高效的单例模式。
func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}

	return instance
}

// 基于 once 实现
func InstanceOnce() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})

	return instance
}
