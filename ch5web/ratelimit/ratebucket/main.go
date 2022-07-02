package main

import (
	"fmt"
	"time"
)

func main() {
	var fillInterval = time.Millisecond * 10
	var capacity = 100
	var tokenBucket = make(chan struct{}, capacity)

	// 填充令牌
	fillToken := func() {
		ticker := time.NewTicker(fillInterval)
		for {
			select {
			case <-ticker.C:
				select {
				case tokenBucket <- struct{}{}:
				default:
				}
				fmt.Println("current token cnt:", len(tokenBucket), time.Now())
			}
		}
	}

	// 取令牌
	TakeAvailable := func(block bool) bool {
		var takenResult bool
		if block {
			select {
			case <-tokenBucket:
				takenResult = true
			}
		} else {
			select {
			case <-tokenBucket:
				takenResult = true
			default:
				takenResult = false
			}
		}

		return takenResult
	}

	go fillToken()
	go func() {
		for i := 0; i < 200; i++ {
			TakeAvailable(true)
		}
	}()
	time.Sleep(time.Second * 3)
}
