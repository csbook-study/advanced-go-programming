package concurrentmemory

import (
	"fmt"
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

func addMutex() {
	var wg sync.WaitGroup
	wg.Add(2)
	go workerm(&wg)
	go workerm(&wg)
	wg.Wait()

	fmt.Println(totalm.value)
}
