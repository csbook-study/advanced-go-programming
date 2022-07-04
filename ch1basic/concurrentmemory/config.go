package concurrentmemory

import (
	"sync/atomic"
	"time"
)

type BaseConfig struct {
}

// 可以继续封装，类型的方法实现（单例模式）
var config atomic.Value

func loadConfig() *BaseConfig {
	return &BaseConfig{}
}

func InitConfig() {
	// 初始化配置信息
	config.Store(loadConfig())

	// 启动一个后台线程，加载更新后的配置信息
	go func() {
		for {
			time.Sleep(time.Minute)
			config.Store(loadConfig())
		}
	}()
}
