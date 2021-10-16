package funcmethodifc

import (
	"fmt"
	"sync"
)

// func
// 具名函数
func Add(a, b int) int {
	return a + b
}

// 匿名函数
var Add2 = func(a, b int) int {
	return a + b
}

// 闭包
func Inc() (v int) {
	defer func() { v++ }()
	return
}

func FuncBasis() {
	// 闭包问题
	// question
	for i := 0; i < 3; i++ {
		defer func() { fmt.Print(i) }()
	}

	// way1: local variable
	for i := 0; i < 3; i++ {
		i := i
		defer func() { fmt.Print(i) }()
	}

	// way2: pass param to defer func
	for i := 0; i < 3; i++ {
		defer func(i int) { fmt.Print(i) }(i)
	}
}

// method
// 1. File
type File struct {
	fd int
}

// 关闭文件
func (f *File) Close() error {
	// ...
	return nil
}

// 读取文件数据
func (f *File) Read(offset int64, data []byte) int {
	// ...
	return 0
}

// 打开文件
func OpenFile(name string) (f *File, err error) {
	// ...
	return
}

// 2. cache
type Cache struct {
	m map[string]string
	sync.Mutex
}

func (p *Cache) Lookup(key string) string {
	p.Lock()
	defer p.Unlock()

	return p.m[key]
}

func MethodBasis() {
	// 方法操作文件
	// 打开文件对象
	f, _ := OpenFile("data")

	// 读取文件数据
	var data []byte
	f.Read(0, data)

	// 关闭文件
	f.Close()
}

func IfcBasis() {

}
