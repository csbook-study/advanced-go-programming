package concurrentmemory

import (
	"fmt"
	"sync"
	"time"
)

func helloWG() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			fmt.Println("hello, world")
		}()
	}

	wg.Wait()
}

func workerwg(wg *sync.WaitGroup, cancel chan bool) {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-cancel:
			return
		case <-time.After(time.Second):
			return
		}
	}
}

func cancelWG() {
	cancel := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go workerwg(&wg, cancel)
	}

	time.Sleep(time.Millisecond)
	close(cancel)
	wg.Wait()
}
