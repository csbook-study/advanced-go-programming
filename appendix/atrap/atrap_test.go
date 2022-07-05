package atrap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestATrap(t *testing.T) {
	t.Run("可变参数是空接口类型", func(t *testing.T) {
		var a = []interface{}{1, 2, 3}

		fmt.Println(a)
		fmt.Println(a...)
	})
	t.Run("数组是值传递", func(t *testing.T) {
		x := [3]int{1, 2, 3}

		func(arr [3]int) {
			arr[0] = 7
			fmt.Println(arr)
		}(x)

		fmt.Println(x)
	})
	t.Run("map遍历是顺序不固定", func(t *testing.T) {
		m := map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
		}

		for k, v := range m {
			println(k, v)
		}
	})
}

func TestATrapFunc(t *testing.T) {
	t.Run("返回值被屏蔽", func(t *testing.T) {
		assert.NotNil(t, errorReturn())
	})
	t.Run("recover必须在defer函数中运行", func(t *testing.T) {
		panicRecover()
	})
	t.Run("通过Sleep来回避并发中的问题：main函数提前退出", func(t *testing.T) {
		sleepWait()
		schedWait()
	})
	t.Run("独占CPU导致其它Goroutine饿死", func(t *testing.T) {
		scheGoroutine()
		selectGoroutine()
	})
	t.Run("不同Goroutine之间不满足顺序一致性内存模型", func(t *testing.T) {
		syncPrint()
	})
	t.Run("闭包错误引用同一个变量", func(t *testing.T) {
		closureVariable()
	})
	t.Run("在循环内部执行defer语句", func(t *testing.T) {
		cycleDefer()
	})
	t.Run("切片会导致整个底层数组被锁定", func(t *testing.T) {
		sliceCopy()
	})
	t.Run("内存地址会变化", func(t *testing.T) {
		memoryAddress()
	})
	t.Run("Goroutine泄露", func(t *testing.T) {
		cancelGoroutine()
	})
}
