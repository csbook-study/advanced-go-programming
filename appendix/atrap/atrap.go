package atrap

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"
)

func errorReturn() (err error) {
	Bar := func() error { return errors.New("err") }
	if err := Bar(); err != nil {
		return err
	}
	return
}

func panicRecover() {
	defer func() {
		recover()
	}()
	panic(1)
}

func sleepWait() {
	go println("hello")
	time.Sleep(time.Second)
}

func schedWait() {
	go println("hello")
	runtime.Gosched()
}

func scheGoroutine() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	for {
		runtime.Gosched()
	}
}

func selectGoroutine() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		os.Exit(0)
	}()

	select {}
}

func syncPrint() {
	var msg string
	var done = make(chan bool)

	setup := func() {
		msg = "hello, world"
		done <- true
	}

	go setup()
	<-done
	println(msg)
}

func closureVariable() {
	for i := 0; i < 5; i++ {
		i := i
		defer func() {
			println(i)
		}()
	}
	for i := 0; i < 5; i++ {
		defer func(i int) {
			println(i)
		}(i)
	}
}

func cycleDefer() {
	for i := 0; i < 5; i++ {
		func() {
			f, err := os.Open("atrap.json")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			d, _ := ioutil.ReadAll(f)
			log.Println(string(d))
		}()
	}
}

func sliceCopy() {
	headerMap := make(map[string][]byte)

	for i := 0; i < 5; i++ {
		name := "atrap.json"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		headerMap[name] = append([]byte{}, data[:1]...)
	}

	// do some thing
}

func memoryAddress() {
	var x int = 42
	var p uintptr = uintptr(unsafe.Pointer(&x))

	runtime.GC()
	var px *int = (*int)(unsafe.Pointer(p))
	println(*px)
}

func cancelGoroutine() {
	ctx, cancel := context.WithCancel(context.Background())

	ch := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		go func() {
			for i := 0; ; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch
	}(ctx)

	for v := range ch {
		fmt.Println(v)
		if v == 5 {
			cancel()
			break
		}
	}
}
