package concurrentmemory

import (
	"fmt"
	"sync"
)

var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "hello, world"
	// output: hello, world false false
	// close(done)
	// output: hello, world true true
	done <- true
}

func syncGoroutine() {
	go aGoroutine()
	ret, ok := <-done
	println(msg, ret, ok)
}

func countGoroutine() {
	var wg sync.WaitGroup
	var limit = make(chan struct{}, 3)
	defer close(limit)

	var work = []func(){
		func() { fmt.Println("goroutine 1") },
		func() { fmt.Println("goroutine 2") },
		func() { fmt.Println("goroutine 3") },
		func() { fmt.Println("goroutine 4") },
	}
	wg.Add(len(work))
	for _, w := range work {
		go func(w func()) {
			defer wg.Done()
			limit <- struct{}{}
			w()
			<-limit
		}(w)
	}
	wg.Wait()

}

// 通过带缓存通道的发送和接收规则可以实现最大并发阻塞，封装函数更加优雅
type gate chan bool

func (g gate) enter() { g <- true }

func (g gate) leave() { <-g }
