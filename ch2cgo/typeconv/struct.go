package main

/*
struct A {
	int type; // type 是 Go 语言的关键字
};

struct A2 {
	int type;    // type 是 Go 语言的关键字
	float _type; // 将屏蔽 CGO 对 type 成员的访问
};
*/
import "C"
import (
	"fmt"
)

func main() {
	var a C.struct_A
	fmt.Println(a._type) // _type 对应 type

	var a2 C.struct_A2
	fmt.Println(a2._type) // _type 对应 _type
}
