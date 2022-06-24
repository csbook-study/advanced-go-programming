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

## 4.5 gRPC 进阶

### 证书认证

gRPC建立在HTTP/2协议之上，对TLS提供了很好的支持。没有启用证书的gRPC服务在和客户端进行的是明文通讯，信息面临被任何第三方监听的风险。为了保障gRPC通信不被第三方监听篡改或伪造，我们可以对服务器启动TLS加密特性。

**公钥认证**

为服务器和客户端分别生成私钥和证书：

```bash
# server 公钥生成
openssl genrsa -out server.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" \
    -key server.key -out server.crt

# client 公钥生成
openssl genrsa -out client.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" \
    -key client.key -out client.crt
```

这种方式，需要提前将服务器的证书告知客户端，这样客户端在链接服务器时才能进行对服务器证书认证。在复杂的网络环境中，服务器证书的传输本身也是一个非常危险的问题。如果在中间某个环节，服务器证书被监听或替换那么对服务器的认证也将不再可靠。

**签名认证**

为了避免证书的传递过程中被篡改，可以通过一个安全可靠的根证书分别对服务器和客户端的证书进行签名。这样客户端或服务器在收到对方的证书后可以通过根证书进行验证证书的有效性。

根证书的生成方式和自签名证书的生成方式类似：

```bash
# ca 生成
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=gobook/CN=github.com" \
    -key ca.key -out ca.crt

# server ca 签名
openssl req -new \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -key server.key \
    -out server.csr
openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in server.csr \
    -out server.crt

# client ca 签名
openssl req -new \
    -subj "/C=GB/L=China/O=client/CN=client.io" \
    -key client.key \
    -out client.csr
openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in client.csr \
    -out client.crt
```

签名的过程中引入了一个新的以.csr为后缀名的文件，它表示证书签名请求文件。在证书签名完成之后可以删除.csr文件。

创建包含SAN的证书

```bash
# ca 生成
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=gobook/CN=github.com" \
    -key ca.key -out ca.crt
    
openssl genrsa -out server.key 2048
openssl req -new -sha256 \
    -key server.key \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -reqexts SAN \
    -config <(cat /etc/pki/tls/openssl.cnf \
        <(printf "\n[SAN]\nsubjectAltName=DNS:server.io,DNS:*.example.com")) \
    -out server.csr
openssl x509 -req -days 3650 \
    -in server.csr -out server.pem \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -extensions SAN \
    -extfile <(cat /etc/pki/tls/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:server.io,DNS:*.example.com"))

openssl genrsa -out client.key 2048
openssl req -new -sha256 \
    -key client.key \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -reqexts SAN \
    -config <(cat /etc/pki/tls/openssl.cnf \
        <(printf "\n[SAN]\nsubjectAltName=DNS:server.io,DNS:*.example.com")) \
    -out client.csr
openssl x509 -req -days 3650 \
    -in client.csr -out client.pem \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -extensions SAN \
    -extfile <(cat /etc/pki/tls/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:server.io,DNS:*.example.com"))
```

### Token认证

gRPC为每个gRPC方法调用提供了认证支持，这样就基于用户Token对不同的方法访问进行权限管理。

要实现对每个gRPC方法进行认证，需要实现grpc.PerRPCCredentials接口：

```go
type PerRPCCredentials interface {
    // GetRequestMetadata gets the current request metadata, refreshing
    // tokens if required. This should be called by the transport layer on
    // each request, and the data should be populated in headers or other
    // context. If a status code is returned, it will be used as the status
    // for the RPC. uri is the URI of the entry point for the request.
    // When supported by the underlying implementation, ctx can be used for
    // timeout and cancellation.
    // TODO(zhaoq): Define the set of the qualified keys instead of leaving
    // it as an arbitrary string.
    GetRequestMetadata(ctx context.Context, uri ...string) (
        map[string]string,    error,
    )
    // RequireTransportSecurity indicates whether the credentials requires
    // transport security.
    RequireTransportSecurity() bool
}
```

在GetRequestMetadata方法中返回认证需要的必要信息。RequireTransportSecurity方法表示是否要求底层使用安全链接。在真实的环境中建议必须要求底层启用安全的链接，否则认证信息有泄露和被篡改的风险。

**详细地认证工作**：首先通过metadata.FromIncomingContext从ctx上下文中获取元信息，然后取出相应的认证信息进行认证。如果认证失败，则返回一个codes.Unauthenticated类型地错误。

### 截取器

gRPC中的grpc.UnaryInterceptor和grpc.StreamInterceptor分别对普通方法和流方法提供了截取器的支持。我们这里简单介绍普通方法的截取器用法。

要实现普通方法的截取器，需要为grpc.UnaryInterceptor的参数实现一个函数：

```go
func filter(ctx context.Context,
    req interface{}, info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    log.Println("fileter:", info)
    return handler(ctx, req)
}
```

函数的ctx和req参数就是每个普通的RPC方法的前两个参数。第三个info参数表示当前是对应的那个gRPC方法，第四个handler参数对应当前的gRPC方法函数。上面的函数中首先是日志输出info参数，然后调用handler对应的gRPC方法函数。

要使用filter截取器函数，只需要在启动gRPC服务时作为参数输入即可：

```go
server := grpc.NewServer(grpc.UnaryInterceptor(filter))
```

如果截取器函数返回了错误，那么该次gRPC方法调用将被视作失败处理。因此，我们可以在截取器中对输入的参数做一些简单的验证工作。同样，也可以对handler返回的结果做一些验证工作。截取器也非常适合前面对Token认证工作。

下面是截取器增加了对gRPC方法异常的捕获：

```go
func filter(
    ctx context.Context, req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    log.Println("fileter:", info)

    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()

    return handler(ctx, req)
}
```

不过gRPC框架中只能为每个服务设置一个截取器，因此所有的截取工作只能在一个函数中完成。开源的grpc-ecosystem项目中的go-grpc-middleware包已经基于gRPC对截取器实现了链式截取器的支持。

以下是go-grpc-middleware包中链式截取器的简单用法

```go
import "github.com/grpc-ecosystem/go-grpc-middleware"

myServer := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        filter1, filter2, ...
    )),
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        filter1, filter2, ...
    )),
)
```

### 和Web服务共存

gRPC构建在HTTP/2协议之上，因此我们可以将gRPC服务和普通的Web服务架设在同一个端口之上。

对于没有启动TLS协议的服务则需要对HTTP2/2特性做适当的调整：

```go
func main() {
    mux := http.NewServeMux()

    h2Handler := h2c.NewHandler(mux, &http2.Server{})
    server = &http.Server{Addr: ":3999", Handler: h2Handler}
    server.ListenAndServe()
}
```

启用普通的https服务器则非常简单：

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(w, "hello")
    })

    http.ListenAndServeTLS(port, "server.crt", "server.key",
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            mux.ServeHTTP(w, r)
            return
        }),
    )
}
```

而单独启用带证书的gRPC服务也是同样的简单：

```go
func main() {
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }

    grpcServer := grpc.NewServer(grpc.Creds(creds))

    ...
}
```

因为gRPC服务已经实现了ServeHTTP方法，可以直接作为Web路由处理对象。如果将gRPC和Web服务放在一起，会导致gRPC和Web路径的冲突，在处理时我们需要区分两类服务。

通过以下方式生成同时支持Web和gRPC协议的路由处理函数：

```go
func main() {
    ...

    http.ListenAndServeTLS(port, "server.crt", "server.key",
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.ProtoMajor != 2 {
                mux.ServeHTTP(w, r)
                return
            }
            if strings.Contains(
                r.Header.Get("Content-Type"), "application/grpc",
            ) {
                grpcServer.ServeHTTP(w, r) // gRPC Server
                return
            }

            mux.ServeHTTP(w, r)
            return
        }),
    )
}
```

首先gRPC是建立在HTTP/2版本之上，如果HTTP不是HTTP/2协议则必然无法提供gRPC支持。同时，每个gRPC调用请求的Content-Type类型会被标注为"application/grpc"类型。

这样就可以在gRPC端口上同时提供Web服务了。

curl触发http请求：

```bash
curl -k --cert client.pem https://localhost:1234
```

## 4.6 gRPC 和 Protobuf 扩展
