package main

// extern int compare(void* a, void* b);
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/huxiangyu99/advanced-go-programming/ch2cgo/qsort"
)

func main() {
	SortMain()
	SortV2Main()
	SortSliceMain()
}

//export compare
func compare(a, b unsafe.Pointer) C.int {
	pa, pb := (*C.int)(a), (*C.int)(b)
	return C.int(*pa - *pb)
}

func SortMain() {
	values := []int32{42, 9, 101, 95, 27, 25}
	qsort.Sort(unsafe.Pointer(&values[0]),
		len(values), int(unsafe.Sizeof(values[0])),
		qsort.CompareFunc(C.compare),
	)
	fmt.Println("sort", values)
}

func SortV2Main() {
	values := []int32{42, 9, 101, 95, 27, 25}
	qsort.SortV2(unsafe.Pointer(&values[0]), len(values), int(unsafe.Sizeof(values[0])),
		func(a, b unsafe.Pointer) int {
			pa, pb := (*int32)(a), (*int32)(b)
			return int(*pa - *pb)
		},
	)
	fmt.Println("sort v2", values)
}

func SortSliceMain() {
	values := []int64{42, 9, 101, 95, 27, 25}
	qsort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	fmt.Println("sort slice", values)
}
