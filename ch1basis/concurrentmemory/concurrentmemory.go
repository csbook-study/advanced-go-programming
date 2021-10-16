package concurrentmemory

import (
	"fmt"
	"runtime"
	"sync"
)

func MutexBasis() {
	var wg sync.WaitGroup
	wg.Add(2)
	go workerm(&wg)
	go workerm(&wg)
	wg.Wait()

	fmt.Println(totalm.value)
}

func AtomicBasis() {
	var wg sync.WaitGroup
	wg.Add(2)
	go workera(&wg)
	go workera(&wg)
	wg.Wait()

	fmt.Println(totala)
}

func SingletonBasis() {
	runtime.GOMAXPROCS(10)

	for i := 0; i < 10; i++ {
		go func() { Instance() }()
	}

	for i := 0; i < 10; i++ {
		go func() { InstanceOnce() }()
	}
}

// 简单的生产者消费者模型：后台生成最新的配置信息；前台多个工作者线程获取最新的配置信息
func ConfigBasis() {
	InitConfig()

	// 处理请求的工作线程始终采用最新的配置信息
	for i := 0; i < 10; i++ {
		go func() {
			// load config
			_ = config.Load().(*BaseConfig)
			// ...
		}()
	}
}

func ChannelBasis() {
	syncGoroutine()
	countGoroutine()
}
