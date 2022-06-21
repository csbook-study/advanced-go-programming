# 第四章 RPC 和 Protobuf

RPC是远程过程调用（Remote Procedure Call）的缩写，即调用远处的一个函数。

## 4.1 RPC 入门

### 跨语言的 RPC

标准库的RPC默认采用Go语言特有的gob编码，因此从其它语言调用Go语言实现的RPC服务将比较困难。在互联网的微服务时代，每个RPC以及服务的使用者都可能采用不同的编程语言，因此跨语言是互联网时代RPC的一个首要条件。得益于RPC的框架设计，Go语言的RPC其实也是很容易实现跨语言支持的。

Go语言的RPC框架有两个比较有特色的设计：一个是RPC数据打包时可以通过插件实现自定义的编码和解码；另一个是RPC建立在抽象的io.ReadWriteCloser接口之上的，我们可以将RPC架设在不同的通讯协议之上。通过官方自带的net/rpc/jsonrpc扩展实现一个跨语言的RPC。

命令行nc模拟请求，查看客户端调用时发送的数据格式和服务端的回包数据：

```bash
# 模拟服务端
nc -l 1234
# {"method":"HelloService.Hello","params":["hello"],"id":0}

# 模拟客户端请求
echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
# {"id":1,"result":"hello, hello","error":null}
```

### HTTP 上的 RPC

Go语言内在的RPC框架已经支持在Http协议上提供RPC服务。但是框架的http服务同样采用了内置的gob协议，并且没有提供采用其它协议的接口，因此从其它语言依然无法访问的。在http协议上提供jsonrpc服务。

模拟一次RPC调用的过程就是向该链接发送一个json字符串：

```bash
# 模拟客户端HTTP请求
curl localhost:1234/jsonrpc -X POST --data '{"method":"HelloService.Hello","params":["hello"],"id":0}'
# {"id":0,"result":"hello, hello","error":null}
```

## 4.2 Protobuf

Protobuf是Protocol Buffers的简称，它是Google公司开发的一种数据描述语言，并于2008年对外开源。Protobuf刚开源时的定位类似于XML、JSON等数据描述语言，通过附带工具生成代码并实现将结构化数据序列化的功能。但是我们更关注的是Protobuf作为接口规范的描述语言，可以作为设计安全的跨语言PRC接口的基础工具。

### Protobuf入门

生成pb桩代码及grpc协议代码：

```bash
protoc --go_out=. hello.proto
protoc --go-grpc_out=. hello.proto
```

### 定制代码生成插件

Protobuf的protoc编译器是通过插件机制实现对不同语言的支持。比如protoc命令出现`--xxx_out`格式的参数，那么protoc将首先查询是否有内置的xxx插件，如果没有内置的xxx插件那么将继续查询当前系统中是否存在protoc-gen-xxx命名的可执行程序，最终通过查询到的插件生成代码。对于Go语言的protoc-gen-go插件来说，里面又实现了一层静态插件系统。比如protoc-gen-go内置了一个gRPC插件，用户可以通过`--go_out=plugins=grpc`参数来生成gRPC相关代码，否则只会针对message生成相关代码。

生成定制化pb桩代码：

```bash
protoc --go-netrpc_out=plugins=netrpc:. hello.proto
```

## 4.3 玩转RPC

### 反向RPC

通常的RPC是基于C/S结构，RPC的服务端对应网络的服务器，RPC的客户端也对应网络客户端。但是对于一些特殊场景，比如在公司内网提供一个RPC服务，但是在外网无法链接到内网的服务器。这种时候我们可以参考类似反向代理的技术，首先从内网主动链接到外网的TCP服务器，然后基于TCP链接向外网提供RPC服务。

反向RPC的内网服务将不再主动提供TCP监听服务，而是首先主动链接到对方的TCP服务器。然后基于每个建立的TCP链接向对方提供RPC服务。

RPC客户端则需要在一个公共的地址提供一个TCP服务，用于接受RPC服务器的链接请求。

### 上下文信息

基于上下文我们可以针对不同客户端提供定制化的RPC服务。我们可以通过为每个链接提供独立的RPC服务来实现对上下文特性的支持。

## 4.4 gRPC 入门

gRPC是Google公司基于Protobuf开发的跨语言的开源RPC框架。gRPC基于HTTP/2协议设计，可以基于一个HTTP/2链接提供多个服务，对于移动设备更加友好。

### gRPC技术栈

![img](https://chai2010.cn/advanced-go-programming-book/images/ch4-1-grpc-go-stack.png)

最底层为TCP或Unix Socket协议，在此之上是HTTP/2协议的实现，然后在HTTP/2协议之上又构建了针对Go语言的gRPC核心库。应用程序通过gRPC插件生产的Stub代码和gRPC核心库通信，也可以直接和gRPC核心库通信。

### gRPC入门

生成pb桩代码及grpc协议代码：

```bash
protoc --go_out=. hello.proto
protoc --go-grpc_out=. hello.proto
protoc --go-grpc_out=. --go_out=. hello.proto
```

gRPC和标准库的RPC框架有一个区别，gRPC生成的接口并不支持异步调用。不过我们可以在多个Goroutine之间安全地共享gRPC底层的HTTP/2链接，因此可以通过在另一个Goroutine阻塞调用的方式模拟异步调用。

### gRPC流

RPC是远程函数调用，因此每次调用的函数参数和返回值不能太大，否则将严重影响每次调用的响应时间。因此传统的RPC方法调用对于上传和下载较大数据量场景并不适合。同时传统RPC模式也不适用于对时间不确定的订阅和发布模式。为此，gRPC框架针对服务器端和客户端分别提供了流特性。

关键字stream指定启用流特性，参数部分是接收客户端参数的流，返回值是返回给客户端的流。

服务端在循环中接收客户端发来的数据，如果遇到io.EOF表示客户端流被关闭，如果函数退出表示服务端流关闭。生成返回的数据通过流发送给客户端，双向流数据的发送和接收都是完全独立的行为。需要注意的是，发送和接收的操作并不需要一一对应，用户可以根据真实场景进行组织代码。
