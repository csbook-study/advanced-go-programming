# 第二章 CGO 编程

## 2.1 快速入门

```go
package main

// #include <stdio.h>
// void SayHello(char* s);
// void SayHelloString(_GoString_ s);
import "C"
import (
  "fmt"
)

func main() {
  C.puts(C.CString("way 1: Hello, World"))
  C.SayHello(C.CString("way 2: Hello, World"))
  C.SayHelloString("way 3: Hello, World")
} 

//export SayHello
func SayHello(s *C.char) {
  fmt.Println(C.GoString(s))
}

//export SayHelloString
func SayHelloString(s string) {
  fmt.Println(int(C._GoStringLen(s)))
  fmt.Println(s)
}
```

## 2.2 CGO 基础

### import "C"

如果在 Go 代码中出现了 import "C" 语句，则表示使用了 CGO 特性，紧临这行语句前面的注释是一种特殊语法，里面包含的是正常的 C 语言代码。当确保 CGO 启用的情况下，还可以在当前目录中包含 C/C++ 对应的源文件。

Go 是强类型语言，所以 CGO 中传递的参数类型必须与声明的类型完全一致，而且传递前必须用 "C" 中的转换函数转换成对应的 C 类型，不能直接传入 Go 中类型的变量。通过虚拟的 C 包导入的 C 语言符号并不需要以大写字母开头，不受 Go 语言的导出规则约束。

CGO 将当前包引用的 C 语言符号都放到了虚拟的 C 包中，同时当前包依赖的其他 Go 语言包内部可能也通过 CGO 引入了相似的虚拟 C 包，但是不同的 Go 语言包引入的虚拟的 C 包之间的类型是不能通用的。

### #cgo

在 import "C" 语句前的注释中可以通过#cgo 语句设置编译阶段和链接阶段的相关参数。编译阶段的参数主要用于定义相关宏和指定头文件检索路径。链接阶段的参数主要是指定库文件检索路径和要链接的库文件。

CFLAGS 部分，-D 部分定义了宏 PNG_DEBUG，值为 1；-I 定义了头文件包含的检索目录（C 头文件检索目录可以是相对路径）。LDFLAGS 部分，-L 指定了链接时库文件检索目录（库文件检索目录需要绝对路径），-l 指定了链接时需要链接 png 库。

```go
// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"
```

\#cgo 语句主要影响 CFLAGS、CPPFLAGS、CXXFLAGS、FFLAGS 和 LDFLAGS 几个编译器环境变量。LDFLAGS 用于设置链接时的参数，除此之外的几个变量用于改变编译阶段的构建参数（CFLAGS 用于针对 C 语言代码设置编译参数）。

对于在 CGO 环境混合使用 C 和 C++ 的用户来说，可能有 3 种不同的编译选项：CFLAGS 对应 C 语言特有的编译选项，CXXFLAGS 对应 C++ 特有的编译选项，CPPFLAGS 则对应 C 和 C++ 共有的编译选项。但是在链接阶段，C 和 C++ 的链接选项是通用的，因此这个时候已经不再有 C 和 C++ 语言的区别，它们的目标文件的类型是相同的。

\#cgo 语句还支持条件选择，当满足某个操作系统或某个 CPU 架构类型时，后面的编译或链接选项生效。可以用 C 语言中常用的技术来处理不同平台之间的差异代码：

```go
/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
  static const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
  static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
  static const char* os = "linux";
#else
#error(unknown os)
#endif
*/
```

### build 标志条件编译

build 标志是在 Go 或 CGO 环境下的 C/C++ 文件开头的一种特殊的注释。条件编译类似于前面通过 #cgo 语句针对不同平台定义的宏，只有在对应平台的宏被定义之后才会构建对应的代码。

```go
// go源文件：
// +build debug

// 编译选项：
go build -tags="debug"
go build -tags="windows debug"
```

## 2.3 类型转换

### 数值类型

在 CGO 中，C 语言的 int 和 long 类型都是对应 4 字节的内存大小，size_t 类型可以当作 Go 语言 uint 无符号整数类型对待。

如果需要在 C 语言中访问 Go 语言的 int 类型，可以通过 GoInt 类型访问，GoInt 类型在 CGO 工具生成的_cgo_export.h 头文件中定义：

```go
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
```

对于比较复杂的 C 语言类型，推荐使用 typedef 关键字提供一个规则的类型命名，这样更利于在 CGO 中访问。

### Go 字符串和切片

在 CGO 生成的_cgo_export.h 头文件中还会为 Go 语言的字符串、切片、字典、接口和通道等特有的数据类型生成对应的 C 语言类型：

```go
typedef struct { const char *p; GoInt n; } GoString;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
```

Go 1.10 针对 Go 字符串增加了一个_GoString_预定义类型，可以降低在 CGO 代码中可能对

_cgo_export.h 头文件产生的循环依赖的风险。

```go
//export helloString
func helloString(s string) {}

extern void helloString(GoString p0);

extern void helloString(_GoString_ p0);
```

因为_GoString_是预定义类型，所以无法通过此类型直接访问字符串的长度和指针等信息。Go 1.10 同时也增加了以下两个函数用于获取字符串结构中的长度和指针信息：

```go
size_t _GoStringLen(_GoString_ s);
const char *_GoStringPtr(_GoString_ s);
```

### 结构体、联合和枚举类型

C 语言的结构体、联合、枚举类型不能作为匿名成员被嵌入到 Go 语言的结构体中。在 Go 语言中，我们可以通过 C.struct_xxx 来访问 C 语言中定义的 struct xxx 结构体类型。结构体的内存布局按照 C 语言的通用对齐规则，C 语言结构体在 32 位 Go 语言环境也按照 32 位对齐规则，在64 位 Go 语言环境按照 64 位对齐规则。对于指定了特殊对齐规则的结构体，无法在 CGO 中访问。

如果结构体的成员名字碰巧是 Go 语言的关键字，则可以通过在成员名开头添加下划线来访问；但是如果有两个成员，一个以 Go 语言关键字命名，另一个刚好是以下划线和 Go 语言关键字命名，那么以 Go 语言关键字命名的成员将无法访问（被屏蔽）。

在 C 语言中，无法直接访问 Go 语言定义的结构体类型。

```go
/*
struct A {
  int type; // type 是 Go 语言的关键字
};

struct B {
  int type;  // type 是 Go 语言的关键字
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

  var b C.struct_B
  fmt.Println(b._type) // _type 对应 _type
}
```

对于联合类型，可以通过 C.union_xxx 来访问 C 语言中定义的 union xxx 类型。但是 Go 语言中并不支持 C 语言联合类型，它们会被转换为对应大小的字节数组。

如果需要操作 C 语言的联合类型变量，一般有 3 种方法：第一种是在 C 语言中定义辅助函数；第二种是通过 Go 语言的"encoding/binary"手工解码成员（需要注意大端小端问题）；第三种是使用 unsafe 包强制转换为对应类型（这是性能最好的方式）。

虽然 unsafe 包访问最简单，性能也最好，但是对于有嵌套联合类型的情况处理会导致问题复杂化。对于复杂的联合类型，推荐通过在 C 语言中定义辅助函数的方式处理。

```go
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
```

对于枚举类型，可以通过 C.enum_xxx 来访问 C 语言中定义的 enum xxx 结构体类型。

在 C 语言中，枚举类型底层对应 int 类型，支持负数类型的值。可以通过 C.ONE、C.TWO 等直接访问定义的枚举值。

```go
/*
enum C {
  ONE,
  TWO,
};
*/
import "C"
import "fmt"

func main() {
  var c C.enum_C = C.TWO
  fmt.Println(c)
  fmt.Println(C.ONE)
  fmt.Println(C.TWO)
}
```

### 数组、字符串和切片

在 C 语言中，数组名其实对应一个指针，指向特定类型特定长度的一段内存，但是这个指针不能被修改。当把数组名传递给一个函数时，实际上传递的是数组第一个元素的地址。C 语言的字符串是一个 char 类型的数组，字符串的长度需要根据表示结尾的 NULL 字符的位置确定。C 语言中没有切片类型。

在 Go 语言中，数组是一种值类型，而且数组的长度是数组类型的一个部分。Go 语言字符串对应一段长度确定的只读 byte 类型的内存。Go 语言的切片则是一个简化版的动态数组。

Go 语言和 C 语言的数组、字符串和切片之间的相互转换可以简化为 Go 语言的切片和 C 语言中指向一定长度内存的指针之间的转换。

```go
func C.CString(string) *C.char
func C.CBytes([]byte) unsafe.Pointer
func C.GoString(*C.char) string
func C.GoStringN(*C.char, C.int) string
func C.GoBytes(unsafe.Pointer, C.int) []byte
```

该组辅助函数都是以克隆的方式运行。当 Go 语言字符串和切片向 C 语言转换时，克隆的内存由 C 语言的 malloc()函数分配，最终可以通过 free()函数释放。当 C 语言字符串或数组向 Go语言转换时，克隆的内存由 Go 语言分配管理。通过该组转换函数，转换前和转换后的内存依然在各自的语言环境中，它们并没有跨越 Go 语言和 C 语言。克隆方式实现转换的优点是接口和内存管理都很简单，缺点是克隆需要分配新的内存和复制操作都会导致额外的开销。

如果字符串或切片对应的底层内存空间由 Go 语言的运行时管理，那么在 C 语言中不能长时间保存 Go 内存对象。

### 指针间的转换

在 C 语言中，不同类型的指针是可以显式或隐式转换的，如果是隐式，则只是会在编译时给出一些警告信息。如果在 Go 语言中两个指针的类型完全一致，则不需要转换可以直接通用。但是 CGO 经常要面对的是两个类型完全不同的指针间的转换，原则上这种操作在纯 Go 语言代码是严格禁止的。

为了实现 X 类型指针到 Y 类型指针的转换，需要借助 unsafe.Pointer 作为中间桥接类型实现不同类型指针之间的转换。unsafe.Pointer 指针类型类似 C 语言中的 void* 类型的指针。

```go
var p *X
var q *Y

q = (*Y)(unsafe.Pointer(p)) // *X => *Y
p = (*X)(unsafe.Pointer(q)) // *Y => *X
```

任何类型的指针都可以通过强制转换为 unsafe.Pointer 指针类型去除原有的类型信息，然后再重新赋予新的指针类型而达到指针间转换的目的。

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-1-x-ptr-to-y-ptr.uml.png)

### 数值和指针的转换

为了严格控制指针的使用，Go 语言禁止将数值类型直接转换为指针类型。不过，Go 语言针对unsafe.Pointer 指针类型特别定义了一个 uintptr 类型。我们可以 uintptr 为中介，实现数值类型到unsafe.Pointer指针类型的转换。

转换分为几个阶段，在每个阶段实现一个小目标：首先是 int32 到 uintptr 类型，然后是 uintptr 到 unsafe.Pointer 指针类型，最后是 unsafe.Pointer 指针类型到 *C.char 类型。

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-2-int32-to-char-ptr.uml.png)

### 切片间的转换

不同切片类型之间转换的思路是先构造一个空的目标切片，然后用原有切片的底层数据填充目标切片。如果类型 X 和 Y 的大小不同，则需要重新设置 Len 和 Cap 属性。需要注意的是，如果 X 或 Y 是空类型，则可能导致除以 0 错误。

```go
var p []X
var q []Y
pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))
pHdr.Data = qHdr.Data
pHdr.Len = qHdr.Len * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
pHdr.Cap = qHdr.Cap * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
```

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-3-x-slice-to-y-slice.uml.png)

## 2.4 函数调用

### Go 调用 C 函数

对于一个启用 CGO 特性的程序，CGO 会构造一个虚拟的 C 包。通过这个虚拟的 C 包可以调用 C 语言函数。

### C 函数的返回值

对于有返回值的 C 函数，我们可以正常获取返回值。

CGO 也针对 <errno.h> 标准库的 errno 宏做了特殊支持：在 CGO 调用 C 函数时如果有两个返回值，那么第二个返回值将对应 errno 错误状态。

### void 函数的返回值

C 语言函数还有一种没有返回值类型的函数，用 void 表示返回值类型。一般情况下，无法获取 void 类型函数的返回值，因为没有返回值可以获取。前面的例子中提到，CGO 对 errno 做了特殊处理，可以通过第二个返回值来获取 C 语言的错误状态。对于 void 类型函数，这个特性依然有效。

我们可以看出 C 语言的 void 类型对应的是当前的 main 包中的_Ctype_void 类型。在 CGO 生成的代码中，_Ctype_void 类型对应一个长度为 0 的数组类型 [0]byte，因此 fmt.Println 输出的是一对表示空数值的方括号。

```go
/*
#include <errno.h>
static int div(int a, int b) {
  if(b == 0) {
    errno = EINVAL;
    return 0;
  }
  return a/b;
}
static void noreturn() {}
*/
import "C"
import "fmt"

func main() {
  v0, err0 := C.div(2, 1)
  fmt.Println(v0, err0)
  v1, err1 := C.div(1, 0)
  fmt.Println(v1, err1)
  v2, err2 := C.noreturn()
  fmt.Printf("%#v, %v, %v", v2, v2, err2)
}
```

### C 调用 Go 导出函数

CGO 生成的_cgo_export.h 文件会包含导出后的 C 语言函数的声明。我们可以在纯 C 源文件中包含_cgo_export.h 文件来引用导出的函数。

当导出 C 语言接口时，需要保证函数的参数和返回值类型都是 C 语言友好的类型，同时返回值不得直接或间接包含 Go 语言内存空间的指针。

## 2.5 内部机制

### CGO 生成的中间文件

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-4-cgo-generated-files.dot.png)

### Go 调用 C 函数

```bash
go tool cgo main.go
```

Go 语言和 C 语言有着不同的内存模型和函数调用规范。其中 _cgo_topofstack() 函数相关的代码用于 C 函数调用后恢复调用栈。_cgo_tsan_acquire() 和 _cgo_tsan_release() 则用于扫描 CGO 相关函数的输入参数和返回值中的指针是否满足规范。runtime.cgocall() 函数是实现 Go 语言到 C 语言函数跨界调用的关键。

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-5-call-c-sum-v1.uml.png)

### C 调用 Go 函数

```bash
go build -buildmode=c-archive -o sum.a sum.go
```

![img](https://chai2010.cn/advanced-go-programming-book/images/ch2-6-call-c-sum-v2.uml.png)

## 2.6 实战：封装 qsort

qsort()快速排序函数是 C 语言的高阶函数，支持用于自定义排序比较函数，可以对任意类型的数组进行排序。

## 2.7 CGO 内存模型

### Go 访问 C 内存

因为 Go 语言实现的限制，无法在 Go 语言中创建大于 2 GB 内存的切片（具体参考 runtime 包中的 makeslice 函数的实现代码）。不过借助 CGO 技术，我们可以在 C 语言环境创建大于 2 GB 的内存，然后转为 Go 语言的切片使用。

因为 C 语言内存空间是稳定的，所以基于 C 语言内存构造的切片也是绝对稳定的，不会因为Go 语言栈的变化而被移动。

```go
/*
#include <stdlib.h>

void* makeslice(size_t memsize) {
  return malloc(memsize);
}
*/
import "C"
import "unsafe"

func makeByteSlize(n int) []byte {
  p := C.makeslice(C.size_t(n))
  return ((*[1<<32 + 1]byte)(p))[0:n:n]
}

func freeByteSlice(p []byte) {
  C.free(unsafe.Pointer(&p[0]))
}

func main() {
  s := makeByteSlize(1<<32 + 1)
  s[len(s)-1] = 12
  print(s[len(s)-1])
  freeByteSlice(s)
}
```

### C 临时访问传入的 Go 内存

借助 C 语言内存稳定的特性，在 C 语言空间先开辟同样大小的内存，然后将 Go 的内存填充到 C 的内存空间，返回的内存也如此处理。

为了简化并高效处理此种向 C 语言传入 Go 语言内存的问题，CGO 针对该场景定义了专门的规则：在 CGO 调用的 C 语言函数返回前，CGO 保证传入的 Go 语言内存在此期间不会发生移动，C 语言函数可以大胆地使用 Go 语言的内存。

```go
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
```

### C 长期持有 Go 指针对象

如果需要在 C 语言中访问 Go 语言内存对象，可以将 Go 语言内存对象在 Go 语言空间映射为一个 int 类型的 ID，然后通过此 ID 来间接访问和控制 Go 语言对象。

### 导出 C 函数不能返回 Go 内存

CGO 默认对返回结果的指针的检查是有代价的，特别是当 CGO 函数返回的结果是一个复杂的数据结构时将花费更多的时间。如果已经确保了 CGO 函数返回的结果是安全的话，那么可以通过设置环境变量 GODEBUG=cgocheck=0 来关闭指针检查行为。

```bash
GODEBUG=cgocheck=0 go run main.go
```

## 2.8 C++类包装

## 2.9 静态库和动态库

## 2.10 编译和链接参数
