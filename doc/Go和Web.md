# 第五章 Go 和 Web

## 5.1 Web 开发简介

Go的Web框架大致可以分为这么两类：

1. Router框架
2. MVC类框架

Go的`net/http`包提供的就是这样的基础功能，写一个简单的`http echo server`只需要30s默认的`net/http`包中的`mux`不支持带参数的路由。

简单地来说，只要你的路由带有参数，并且这个项目的API数目超过了10，就尽量不要使用`net/http`中默认的路由。在Go开源界应用最广泛的router是httpRouter，很多开源的router框架都是基于httpRouter进行一定程度的改造的成果。

开源界有这么几种框架，第一种是对httpRouter进行简单的封装，然后提供定制的中间件和一些简单的小工具集成比如gin，主打轻量，易学，高性能。第二种是借鉴其它语言的编程风格的一些MVC类框架，例如beego，方便从其它语言迁移过来的程序员快速上手，快速开发。还有一些框架功能更为强大，除了数据库schema设计，大部分代码直接生成，例如goa。

## 5.2 router 请求路由

在常见的Web框架中，router是必备的组件。Go语言圈子里router也时常被称为`http`的multiplexer。

RESTful是几年前刮起的API设计风潮，在RESTful中除了GET和POST之外，还使用了HTTP协议定义的几种其它的标准化语义。具体包括：

```go
const (
    MethodGet     = "GET"
    MethodHead    = "HEAD"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodPatch   = "PATCH" // RFC 5789
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
)
```

来看看RESTful中常见的请求路径：

```bash
GET /repos/:owner/:repo/comments/:id/reactions

POST /projects/:project_id/columns

PUT /user/starred/:owner/:repo

DELETE /user/starred/:owner/:repo
```

RESTful风格的API重度依赖请求路径。会将很多参数放在请求URI中。除此之外还会使用很多并不那么常见的HTTP状态码。

### httprouter

较流行的开源go Web框架大多使用httprouter，或是基于httprouter的变种对路由进行支持。

因为httprouter中使用的是显式匹配，所以在设计路由的时候需要规避一些会导致路由冲突的情况，例如：

```bash
conflict:
GET /user/info/:name
GET /user/:id

no conflict:
GET /user/info/:name
POST /user/:id
```

简单来讲的话，如果两个路由拥有一致的http方法(指 GET/POST/PUT/DELETE)和请求路径前缀，且在某个位置出现了A路由是wildcard（指:id这种形式）参数，B路由则是普通字符串，那么就会发生路由冲突。路由冲突会在初始化阶段直接panic。

```bash
# panic: wildcard route ':id' conflicts with existing children in path '/user/:id'
```

还有一点需要注意，因为httprouter考虑到字典树的深度，在初始化时会对参数的数量进行限制，所以在路由中的参数数目不能超过255，否则会导致httprouter无法识别后续的参数。不过这一点上也不用考虑太多，毕竟URI是人设计且给人来看的，相信没有长得夸张的URI能在一条路径中带有200个以上的参数。

除支持路径中的wildcard参数之外，httprouter还可以支持`*`号来进行通配，不过`*`号开头的参数只能放在路由的结尾，例如下面这样：

```bash
Pattern: /src/*filepath

 /src/                     filepath = ""
 /src/somefile.go          filepath = "somefile.go"
 /src/subdir/somefile.go   filepath = "subdir/somefile.go"
```

这种设计在RESTful中可能不太常见，主要是为了能够使用httprouter来做简单的HTTP静态文件服务器。

除了正常情况下的路由支持，httprouter也支持对一些特殊情况下的回调函数进行定制，例如404的时候：

```go
r := httprouter.New()
r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("oh no, not found"))
})
```

或者内部panic的时候：

```go
r.PanicHandler = func(w http.ResponseWriter, r *http.Request, c interface{}) {
    log.Printf("Recovering from panic, Reason: %#v", c.(error))
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(c.(error).Error()))
}
```

目前开源界最为流行（star数最多）的Web框架[gin](https://github.com/gin-gonic/gin)使用的就是httprouter的变种。

### 原理

httprouter和众多衍生router使用的数据结构被称为压缩字典树（Radix Tree）。

典型的字典树结构：

![trie tree](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-trie.png)

字典树常用来进行字符串检索，例如用给定的字符串序列建立字典树。**对于目标字符串，只要从根节点开始深度优先搜索，即可判断出该字符串是否曾经出现过，时间复杂度为`O(n)`，n可以认为是目标字符串的长度。**为什么要这样做？字符串本身不像数值类型可以进行数值比较，两个字符串对比的时间复杂度取决于字符串长度。如果不用字典树来完成上述功能，要对历史字符串进行排序，再利用二分查找之类的算法去搜索，时间复杂度只高不低。可认为字典树是一种空间换时间的典型做法。

**普通的字典树有一个比较明显的缺点，就是每个字母都需要建立一个孩子节点，这样会导致字典树的层数比较深，压缩字典树相对好地平衡了字典树的优点和缺点。**是典型的压缩字典树结构：

![radix tree](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-radix.png)

每个节点上不只存储一个字母了，这也是压缩字典树中“压缩”的主要含义。使用压缩字典树可以减少树的层数，同时因为每个节点上数据存储也比通常的字典树要多，所以程序的局部性较好（一个节点的path加载到cache即可进行多个字符的对比)，从而对CPU缓存友好。

### 压缩字典树创建过程

跟踪httprouter中一个典型的压缩字典树的创建过程，路由设定如下：

```bash
PUT /user/installations/:installation_id/repositories/:repository_id

GET /marketplace_listing/plans/
GET /marketplace_listing/plans/:id/accounts
GET /search
GET /status
GET /support
```

#### root节点创建

httprouter的Router结构体中存储压缩字典树使用的是下述数据结构：

```go
// 略去了其它部分的 Router struct
type Router struct {
    // ...
    trees map[string]*node
    // ...
}
```

`trees`中的`key`即为HTTP 1.1的RFC中定义的各种方法`GET/HEAD/OPTIONS/POST/PUT/PATCH/DELETE`。

**每一种方法对应的都是一棵独立的压缩字典树，这些树彼此之间不共享数据。**具体到我们上面用到的路由，`PUT`和`GET`是两棵树而非一棵。

radix的节点类型为`*httprouter.node`，为了说明方便，我们留下了目前关心的几个字段：

```bash
path: 当前节点对应的路径中的字符串

wildChild: 子节点是否为参数节点，即 wildcard node，或者说 :id 这种类型的节点

nType: 当前节点类型，有四个枚举值: 分别为 static/root/param/catchAll。
    static                   // 非根节点的普通字符串节点
    root                     // 根节点
    param                    // 参数节点，例如 :id
    catchAll                 // 通配符节点，例如 *anyway

indices：子节点索引，当子节点为非参数类型，即本节点的wildChild为false时，会将每个子节点的首字母放在该索引数组。说是数组，实际上是个string。
```

#### 子节点插入

插入`GET /marketplace_listing/plans`时，结构如下：

![get radix step 1](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-radix-get-1.png)

因为第一个路由没有参数，path都被存储到根节点上了。所以只有一个节点。

然后插入`GET /marketplace_listing/plans/:id/accounts`，新的路径与之前的路径有共同的前缀，且可以直接在之前叶子节点后进行插入，那么结果也很简单，插入后的树结构如下：

![get radix step 2](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-radix-get-2.png)

由于`:id`这个节点只有一个字符串的普通子节点，所以indices还依然不需要处理。

上面这种情况比较简单，新的路由可以直接作为原路由的子节点进行插入。

#### 边分裂

接下来我们插入`GET /search`，这时会导致树的边分裂，插入后的树结构如下：

![get radix step 3](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-radix-get-3.png)

原有路径和新的路径在初始的`/`位置发生分裂，这样需要把原有的root节点内容下移，再将新路由 `search`同样作为子节点挂在root节点之下。这时候因为子节点出现多个，root节点的indices提供子节点索引，这时候该字段就需要派上用场了。"ms"代表子节点的首字母分别为m（marketplace）和s（search)。

我们一口作气，把`GET /status`和`GET /support`也插入到树中。这时候会导致在`search`节点上再次发生分裂，最终结果如下：

![get radix step 4](https://chai2010.cn/advanced-go-programming-book/images/ch6-02-radix-get-4.png)

#### 子节点冲突处理

在路由本身只有字符串的情况下，不会发生任何冲突。只有当路由中含有wildcard（类似 :id）或者catchAll的情况下才可能冲突。这一点在前面已经提到了。

子节点的冲突处理很简单，分几种情况：

1. 在插入wildcard节点时，父节点的children数组非空且wildChild被设置为false。例如：`GET /user/getAll`和`GET /user/:id/getAddr`，或者`GET /user/*aaa`和`GET /user/:id`。
2. 在插入wildcard节点时，父节点的children数组非空且wildChild被设置为true，但该父节点的wildcard子节点要插入的wildcard名字不一样。例如：`GET /user/:id/info`和`GET /user/:name/info`。
3. 在插入catchAll节点时，父节点的children非空。例如：`GET /src/abc`和`GET /src/*filename`，或者`GET /src/:id`和`GET /src/*filename`。
4. 在插入static节点时，父节点的wildChild字段被设置为true。
5. 在插入static节点时，父节点的children非空，且子节点nType为catchAll。

只要发生冲突，都会在初始化的时候panic。例如，在插入我们臆想的路由`GET /marketplace_listing/plans/ohyes`时，出现第4种冲突情况：它的父节点`marketplace_listing/plans/`的wildChild字段为true。

## 5.3 中间件

Web框架中的中间件(middleware)技术原理。

### 使用中间件剥离非业务逻辑

对于大多数的场景来讲，非业务的需求都是在http请求处理前做一些事情，并且在响应完成之后做一些事情。

```go
http.Handle("/", timeMiddleware(http.HandlerFunc(hello)))
```

任何方法实现了`ServeHTTP`，即是一个合法的`http.Handler`，读到这里你可能会有一些混乱，我们先来梳理一下http库的`Handler`，`HandlerFunc`和`ServeHTTP`的关系：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

只要你的handler函数签名是：

```go
func (ResponseWriter, *Request)
```

那么这个`handler`和`http.HandlerFunc()`就有了一致的函数签名，可以将该`handler()`函数进行类型转换，转为`http.HandlerFunc`。而`http.HandlerFunc`实现了`http.Handler`这个接口。在`http`库需要调用你的handler函数来处理http请求时，会调用`HandlerFunc()`的`ServeHTTP()`函数，可见一个请求的基本调用链是这样的：

```go
h = getHandler() => h.ServeHTTP(w, r) => h(w, r)
```

上面提到的把自定义`handler`转换为`http.HandlerFunc()`这个过程是必须的，因为我们的`handler`没有直接实现`ServeHTTP`这个接口。上面的代码中我们看到的HandleFunc(注意HandlerFunc和HandleFunc的区别)里也可以看到这个强制转换过程：

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}

// 调用

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    mux.Handle(pattern, HandlerFunc(handler))
}
```

中间件要做的事情就是通过一个或多个函数对handler进行包装，返回一个包括了各个中间件逻辑的函数链。

```go
customizedHandler = logger(timeout(ratelimit(helloHandler)))
```

这个函数链在执行过程中的上下文：

![img](https://chai2010.cn/advanced-go-programming-book/images/ch6-03-middleware_flow.png)

这个流程在进行请求处理的时候就是不断地进行函数压栈再出栈，有一些类似于递归的执行流：

```bash
[exec of logger logic]           函数栈: []

[exec of timeout logic]          函数栈: [logger]

[exec of ratelimit logic]        函数栈: [timeout/logger]

[exec of helloHandler logic]     函数栈: [ratelimit/timeout/logger]

[exec of ratelimit logic part2]  函数栈: [timeout/logger]

[exec of timeout logic part2]    函数栈: [logger]

[exec of logger logic part2]     函数栈: []
```

函数套函数的用法不是很美观，同时也不具备什么可读性。

### 更优雅的中间件写法

```go
r = NewRouter()
r.Use(logger)
r.Use(timeout)
r.Use(ratelimit)
r.Add("/", helloHandler)
```

通过多步设置，我们拥有了和上一节差不多的执行函数链。胜在直观易懂，如果我们要增加或者删除中间件，只要简单地增加删除对应的`Use()`调用就可以了。非常方便。

### 适合中间件做的事情

以较流行的开源Go语言框架chi为例：

```bash
compress.go
  => 对http的响应体进行压缩处理
heartbeat.go
  => 设置一个特殊的路由，例如/ping，/healthcheck，用来给负载均衡一类的前置服务进行探活
logger.go
  => 打印请求处理处理日志，例如请求处理时间，请求路由
profiler.go
  => 挂载pprof需要的路由，如`/pprof`、`/pprof/trace`到系统中
realip.go
  => 从请求头中读取X-Forwarded-For和X-Real-IP，将http.Request中的RemoteAddr修改为得到的RealIP
requestid.go
  => 为本次请求生成单独的requestid，可一路透传，用来生成分布式调用链路，也可用于在日志中串连单次请求的所有逻辑
timeout.go
  => 用context.Timeout设置超时时间，并将其通过http.Request一路透传下去
throttler.go
  => 通过定长大小的channel存储token，并通过这些token对接口进行限流
```

## 5.4 请求校验
