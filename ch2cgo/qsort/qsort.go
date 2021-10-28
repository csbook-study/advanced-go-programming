package qsort

/*
#include <stdlib.h>

typedef int (*qsort_cmp_func_t)(const void* a, const void* b);
extern int _cgo_qsort_compare(void* a, void* b);
extern int _cgo_qsort_compare_v2(void* a, void* b);
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type CompareFunc C.qsort_cmp_func_t

func Sort(base unsafe.Pointer, num, size int, cmp CompareFunc) {
	C.qsort(base, C.size_t(num), C.size_t(size), C.qsort_cmp_func_t(cmp))
}

func SortV2(base unsafe.Pointer, num, size int, cmp func(a, b unsafe.Pointer) int) {
	go_qsort_compare_info.Lock()
	defer go_qsort_compare_info.Unlock()
	go_qsort_compare_info.fn = cmp
	C.qsort(base, C.size_t(num), C.size_t(size),
		C.qsort_cmp_func_t(C._cgo_qsort_compare),
	)
}

func Slice(slice interface{}, less func(a, b int) bool) {
	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("qsort called with non-slice value of type %T", slice))
	}
	if sv.Len() == 0 {
		return
	}
	go_qsort_compare_info_v2.Lock()
	defer go_qsort_compare_info_v2.Unlock()
	defer func() {
		go_qsort_compare_info_v2.base = nil
		go_qsort_compare_info_v2.elemnum = 0
		go_qsort_compare_info_v2.elemsize = 0
		go_qsort_compare_info_v2.less = nil
	}()
	go_qsort_compare_info_v2.base = unsafe.Pointer(sv.Index(0).Addr().Pointer())
	go_qsort_compare_info_v2.elemnum = sv.Len()
	go_qsort_compare_info_v2.elemsize = int(sv.Type().Elem().Size())
	go_qsort_compare_info_v2.less = less
	C.qsort(
		go_qsort_compare_info_v2.base,
		C.size_t(go_qsort_compare_info_v2.elemnum),
		C.size_t(go_qsort_compare_info_v2.elemsize),
		C.qsort_cmp_func_t(C._cgo_qsort_compare_v2),
	)
}
