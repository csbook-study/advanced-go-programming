package concurrentmode

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func ProducerConsumerBasis() {
	ch := make(chan int, 64)

	go Producer(3, ch)
	go Producer(5, ch)
	go Consumer(ch)

	// wait to process
	time.Sleep(time.Millisecond)

	// ctrl+c exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit %v\n", <-sig)
}

func PublishSubscribeBasis() {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello, world")
	p.Publish("hello, golang")

	go func() {
		for msg := range all {
			fmt.Println("all: ", msg)
		}
	}()
	go func() {
		for msg := range golang {
			fmt.Println("golang: ", msg)
		}
	}()

	// wait to process
	time.Sleep(time.Second)
}

// 素数筛
func PrimeSieveBasis() {
	ch := GeneralNatural()
	for i := 0; i < 100; i++ {
		prime := <-ch
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ch, prime)
	}
}

// 素数筛
func PrimeSieveContextBasis() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	ch := GeneralNaturalContext(ctx)
	for i := 0; i < 100; i++ {
		prime := <-ch
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilterContext(ctx, ch, prime)
	}

	cancel()
}
