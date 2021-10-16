package concurrentmemory

import (
	"sync"
	"sync/atomic"
)

var totala uint64

func workera(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&totala, i)
	}
}
