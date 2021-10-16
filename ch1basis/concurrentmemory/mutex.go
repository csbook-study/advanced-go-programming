package concurrentmemory

import (
	"sync"
)

var totalm struct {
	sync.Mutex
	value int
}

func workerm(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 100; i++ {
		totalm.Lock()
		totalm.value++
		totalm.Unlock()
	}
}
