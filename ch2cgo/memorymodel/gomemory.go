package main

/*
#include <stdio.h>
#include <stdlib.h>

void printString(const char* s) {
	printf("%s\n", s);
}

void printStringSafe(const char* s, int n) {
	int i;
	for(i = 0; i < n; i++) {
		putchar(s[i]);
	}
	putchar('\n');
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

func main() {
	s := "hello"
	printString(s)
	printStringSafe(s)
}

func printString(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.printString(cs)
}

func printStringSafe(s string) {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	C.printStringSafe((*C.char)(unsafe.Pointer(p.Data)), C.int(len(s)))
}
