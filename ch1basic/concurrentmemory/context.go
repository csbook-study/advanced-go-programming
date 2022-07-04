package concurrentmemory

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func workerc(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func cancelContext() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go workerc(ctx, &wg)
	}

	time.Sleep(time.Millisecond)
	cancel()
	wg.Wait()
}
