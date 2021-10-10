package arraystringslice

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

func ArrayBasis() {
	var a [3]int                    // 定义长度为 3 的 int 型数组，元素为 0、0、0
	var b = [...]int{1, 2, 3}       // 定义长度为 3 的 int 型数组，元素为 1、2、3
	var c = [...]int{2: 3, 1: 2}    // 定义长度为 3 的 int 型数组，元素为 0、2、3
	var d = [...]int{1, 2, 4: 5, 6} // 定义长度为 6 的 int 型数组，元素为 1、2、0、0、5、6

	// output: [0 0 0] [1 2 3] [0 2 3] [1 2 0 0 5 6]
	fmt.Println(a, b, c, d)
}

func StringBasis() {
	s := "hello, world"                              // Len 12
	hello, world := s[:5], s[7:]                     // Len 5 5
	s1, s2 := "hello, world"[:5], "hello, world"[7:] // Len 5 5

	// output: len(s): 12 data: 0xc000088220
	fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len,
		"data:", &(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	// output: len(hello): 5 data: 0xc000088230
	fmt.Println("len(hello):", (*reflect.StringHeader)(unsafe.Pointer(&hello)).Len,
		"data:", &(*reflect.StringHeader)(unsafe.Pointer(&hello)).Data)
	// output: len(world): 5 data: 0xc000088240
	fmt.Println("len(world):", (*reflect.StringHeader)(unsafe.Pointer(&world)).Len,
		"data:", &(*reflect.StringHeader)(unsafe.Pointer(&world)).Data)
	// output: len(s1): 5 data: 0xc000088250
	fmt.Println("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len,
		"data:", &(*reflect.StringHeader)(unsafe.Pointer(&s1)).Data)
	// output: len(s2): 5 data: 0xc000088260
	fmt.Println("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len,
		"data:", &(*reflect.StringHeader)(unsafe.Pointer(&s2)).Data)
}

func SliceBasis() {
	// 1. define
	var a []int               // nil 切片，和 nil 相等，一般用来表示一个不存在的切片
	var b = []int{}           // 空切片，和 nil 不相等，一般用来表示一个空集合
	var c = []int{1, 2, 3}    // 3个元素的切片，len 和 cap 都为 3
	var d = c[1:2]            // 2个元素的切片，len 为 1，cap 为 2
	var e = c[0:2:cap(c)]     // 2个元素的切片，len 为 2，cap 为 3
	var f = c[:0]             // 0个元素的切片，len 为 0，cap 为 3
	var g = make([]int, 3)    // 3个元素的切片，len 和 cap 都为 3
	var h = make([]int, 2, 3) // 2个元素的切片，len 为 2，cap 为 3
	var i = make([]int, 0, 3) // 0个元素的切片，len 为 0，cap 为 3

	// output: [] [] [1 2 3] 3 [2] 2 [1 2] 3 [] 3 [0 0 0] 3 [0 0] 3 [] 3
	fmt.Println(a, b, c, cap(c), d, cap(d), e, cap(e), f, cap(f), g, cap(g), h, cap(h), i, cap(i))

	// 2. append
	index, value := 0, 99
	// 向切片 index 位置插入一个元素 value
	a = append(a, 0)             // 切片扩展一个空间
	copy(a[index+1:], a[index:]) // a[index:] 向后移动一个元素
	a[index] = value             // 设置新添加的元素

	// output: [99] 1
	fmt.Println(a, cap(a))

	values := []int{23, 24, 25, 0, 0, 26, 27, 28}
	// 向切片 index 位置插入多个元素 values
	a = append(a, values...)               // 切片扩展一个空间
	copy(a[index+len(values):], a[index:]) // a[index:] 向后移动 len(values) 个元素
	copy(a[index:], values)                // 设置新添加的切片

	// output: [23 24 25 0 0 26 27 28 99] 10
	fmt.Println(a, cap(a))

	// 3. delete
	index, n := 1, 2
	// 删除切片中间 index 位置的元素（n 个）
	a = append(a[:index], a[index+1:]...)
	a = append(a[:index], a[index+n:]...)
	a = a[:index+copy(a[index:], a[index+1:])]
	a = a[:index+copy(a[index:], a[index+n:])]

	// output: [23 28 99] 10
	fmt.Println(a, cap(a))

	// 4. technique
	// TrimSlice Filter
	s := []byte{'1', '2', '3', ' ', ' ', '4'}
	s = TrimSlice(s)

	// output: 1234
	fmt.Println(string(s))

	// 5. memory

	// 6. type conversion
	// SortFloat64FastV1 SortFloat64FastV2
	nums := []float64{4, 2, 5, 7, 1, 9, 0}
	SortFloat64FastV1(nums)
	SortFloat64FastV2(nums)

	// output: [0 1 2 4 5 7 9] 7
	fmt.Println(nums, cap(nums))
}

func TrimSlice(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

func Filter(s []byte, fn func(x byte) bool) []byte {
	b := s[:0]
	for _, x := range s {
		if !fn(x) {
			b = append(b, x)
		}
	}
	return b
}

func SortFloat64FastV1(a []float64) {
	// 强制类型转换
	var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]

	// 以 int 方式给 float64 排序
	sort.Ints(b)
}

func SortFloat64FastV2(a []float64) {
	// 通过 reflect.SliceHeader 更新切片头部信息实现转换
	var c []int
	aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*cHdr = *aHdr

	// 以 int 方式给 float64 排序
	sort.Ints(c)
}
