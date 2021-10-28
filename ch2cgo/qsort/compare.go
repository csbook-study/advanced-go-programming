package qsort

import "C"
import (
	"sync"
	"unsafe"
)

var go_qsort_compare_info struct {
	fn func(a, b unsafe.Pointer) int
	sync.Mutex
}

//export _cgo_qsort_compare
func _cgo_qsort_compare(a, b unsafe.Pointer) C.int {
	return C.int(go_qsort_compare_info.fn(a, b))
}

var go_qsort_compare_info_v2 struct {
	base     unsafe.Pointer
	elemnum  int
	elemsize int
	less     func(a, b int) bool
	sync.Mutex
}

//export _cgo_qsort_compare_v2
func _cgo_qsort_compare_v2(a, b unsafe.Pointer) C.int {
	var (
		// array memory is locked
		base     = uintptr(go_qsort_compare_info_v2.base)
		elemsize = uintptr(go_qsort_compare_info_v2.elemsize)
	)
	i := int((uintptr(a) - base) / elemsize)
	j := int((uintptr(b) - base) / elemsize)
	switch {
	case go_qsort_compare_info_v2.less(i, j): // v[i] < v[j]
		return -1
	case go_qsort_compare_info_v2.less(j, i): // v[i] > v[j]
		return +1
	default:
		return 0
	}
}
