package main

// #include <stdio.h>
// #include "hello.h"
// void SayHello(char* s);
// void SayHelloString(_GoString_ s);
import "C"
import (
	"fmt"
)

func main() {
	C.puts(C.CString("way 1: Hello, World"))
	C.SayHi(C.CString("way 2: Hello, World"))
	C.SayHello(C.CString("way 3: Hello, World"))
	C.SayHelloString("way 4: Hello, World")
}

//export SayHello
func SayHello(s *C.char) {
	fmt.Println(C.GoString(s))
}

//export SayHelloString
func SayHelloString(s string) {
	fmt.Print(s)
	fmt.Println(int(C._GoStringLen(s)))
}
