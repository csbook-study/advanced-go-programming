package concurrentmemory

import (
	"fmt"
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

func addAtomic() {
	var wg sync.WaitGroup
	wg.Add(2)
	go workera(&wg)
	go workera(&wg)
	wg.Wait()

	fmt.Println(totala)
}
