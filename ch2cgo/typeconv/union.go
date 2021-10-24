package main

/*
union B {
	int i;
	float f;
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var b C.union_B
	fmt.Printf("%T\n", b) // [4]uint8
	fmt.Println("b.i:", *(*C.int)(unsafe.Pointer(&b)))
	fmt.Println("b.f:", *(*C.float)(unsafe.Pointer(&b)))
}
