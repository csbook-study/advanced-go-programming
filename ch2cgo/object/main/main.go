package main

/*
extern char* NewGoString(char* );
extern void FreeGoString(char* );
extern void PrintGoString(char* );
static void printString(char* s) {
	char* gs = NewGoString(s);
	PrintGoString(gs);
	FreeGoString(gs);
}
*/
import "C"
import (
	"unsafe"

	"github.com/huxiangyu99/advanced-go-programming/ch2cgo/object"
)

//export NewGoString
func NewGoString(s *C.char) *C.char {
	gs := C.GoString(s)
	id := object.NewObjectId(gs)
	return (*C.char)(unsafe.Pointer(uintptr(id)))
}

//export FreeGoString
func FreeGoString(p *C.char) {
	id := object.ObjectId(uintptr(unsafe.Pointer(p)))
	id.Free()
}

//export PrintGoString
func PrintGoString(s *C.char) {
	id := object.ObjectId(uintptr(unsafe.Pointer(s)))
	gs := id.Get().(string)
	print(gs)
}

func main() {
	C.printString(C.CString("hello"))
}
