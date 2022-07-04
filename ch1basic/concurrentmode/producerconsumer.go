package concurrentmode

import "fmt"

// Producer 生成 factor 整数倍的序列
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- factor * i
	}
}

// Consumer
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
