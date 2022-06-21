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
