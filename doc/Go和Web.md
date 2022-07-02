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

### validator库原理

定义如下的结构体：

```go
type Nested struct {
    Email string `validate:"email"`
}
type T struct {
    Age    int `validate:"eq=10"`
    Nested Nested
}
```

validator 树结构：

![struct-tree](https://chai2010.cn/advanced-go-programming-book/images/ch6-04-validate-struct-tree.png)

通过反射对结构体遍历。

## 5.5 和数据库打交道

### 从 database/sql 讲起

Go官方提供了`database/sql`包来给用户进行和数据库打交道的工作，`database/sql`库实际只提供了一套操作数据库的接口和规范，例如抽象好的SQL预处理（prepare），连接池管理，数据绑定，事务，错误处理等等。官方并没有提供具体某种数据库实现的协议支持。

和具体的数据库，例如MySQL打交道，还需要再引入MySQL的驱动，像下面这样：

```go
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "user:password@/dbname")
```

```go
import _ "github.com/go-sql-driver/mysql"
```

这条import语句会调用了`mysql`包的`init`函数，做的事情也很简单：

```go
func init() {
    sql.Register("mysql", &MySQLDriver{})
}
```

在`sql`包的全局`map`里把`mysql`这个名字的`driver`注册上。`Driver`在`sql`包中是一个接口：

```go
type Driver interface {
    Open(name string) (Conn, error)
}
```

调用`sql.Open()`返回的`db`对象就是这里的`Conn`。

```go
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```

### 提高生产效率的ORM和SQL Builder

对象关系映射（英语：Object Relational Mapping，简称ORM，或O/RM，或O/R mapping），是一种程序设计技术，用于实现面向对象编程语言里不同类型系统的数据之间的转换。从效果上说，它其实是创建了一个可在编程语言里使用的“虚拟对象数据库”。

相比ORM来说，SQL Builder在SQL和项目可维护性之间取得了比较好的平衡。首先sql builder不像ORM那样屏蔽了过多的细节，其次从开发的角度来讲，SQL Builder进行简单封装后也可以非常高效地完成开发，举个例子：

```go
where := map[string]interface{} {
    "order_id > ?" : 0,
    "customer_id != ?" : 0,
}
limit := []int{0,100}
orderBy := []string{"id asc", "create_time desc"}

orders := orderModel.GetList(where, limit, orderBy)
```

### 脆弱的数据库

无论是ORM还是SQL Builder都有一个致命的缺点，就是没有办法进行系统上线的事前sql审核。虽然很多ORM和SQL Builder也提供了运行期打印sql的功能，但只在查询的时候才能进行输出。而SQL Builder和ORM本身提供的功能太过灵活。使得你不可能通过测试枚举出所有可能在线上执行的sql。

## 5.6 服务流量限制

计算机程序可依据其瓶颈分为磁盘IO瓶颈型，CPU计算瓶颈型，网络带宽瓶颈型，分布式场景下有时候也会外部系统而导致自身瓶颈。

Web系统打交道最多的是网络，无论是接收，解析用户请求，访问存储，还是把响应数据返回给用户，都是要走网络的。在没有`epoll/kqueue`之类的系统提供的IO多路复用接口之前，多个核心的现代计算机最头痛的是C10k问题，C10k问题会导致计算机没有办法充分利用CPU来处理更多的用户连接，进而没有办法通过优化程序提升CPU利用率来处理更多的请求。

自从Linux实现了`epoll`，FreeBSD实现了`kqueue`，这个问题基本解决了，我们可以借助内核提供的API轻松解决当年的C10k问题，也就是说如今如果你的程序主要是和网络打交道，那么瓶颈一定在用户程序而不在操作系统内核。

wrk对http服务压测，多次测试的结果在4万左右的QPS浮动，响应时间最多也就是40ms左右。压测命令：

```bash
wrk -c 10 -d 10s -t10 http://localhost:9090
# Running 10s test @ http://localhost:9090
#   10 threads and 10 connections
#   Thread Stats   Avg      Stdev     Max   +/- Stdev
#     Latency   334.76us    1.21ms  45.47ms   98.27%
#     Req/Sec     4.42k   633.62     6.90k    71.16%
#   443582 requests in 10.10s, 54.15MB read
# Requests/sec:  43911.68
# Transfer/sec:      5.36MB
```

对于IO/Network瓶颈类的程序，其表现是网卡/磁盘IO会先于CPU打满，这种情况即使优化CPU的使用也不能提高整个系统的吞吐量，只能提高磁盘的读写速度，增加内存大小，提升网卡的带宽来提升整体性能。而CPU瓶颈类的程序，则是在存储和网卡未打满之前CPU占用率先到达100%，CPU忙于各种计算任务，IO设备相对则较闲。

无论哪种类型的服务，在资源使用到极限的时候都会导致请求堆积，超时，系统hang死，最终伤害到终端用户。对于分布式的Web服务来说，瓶颈还不一定总在系统内部，也有可能在外部。非计算密集型的系统往往会在关系型数据库环节失守，而这时候Web模块本身还远远未达到瓶颈。

不管我们的服务瓶颈在哪里，最终要做的事情都是一样的，那就是流量限制。

### 常见的流量限制手段

流量限制的手段有很多，最常见的：漏桶、令牌桶两种：

1. 漏桶是指我们有一个一直装满了水的桶，每过固定的一段时间即向外漏一滴水。如果你接到了这滴水，那么你就可以继续服务请求，如果没有接到，那么就需要等待下一滴水。
2. 令牌桶则是指匀速向桶中添加令牌，服务请求时需要从桶中获取令牌，令牌的数目可以按照需要消耗的资源进行相应的调整。如果没有令牌，可以选择等待，或者放弃。

这两种方法看起来很像，不过还是有区别的。漏桶流出的速率固定，而令牌桶只要在桶中有令牌，那就可以拿。也就是说令牌桶是允许一定程度的并发的，比如同一个时刻，有100个用户请求，只要令牌桶中有100个令牌，那么这100个请求全都会放过去。令牌桶在桶中没有令牌的情况下也会退化为漏桶模型。

![token bucket](https://chai2010.cn/advanced-go-programming-book/images/ch5-token-bucket.png)

实际应用中令牌桶应用较为广泛，开源界流行的限流器大多数都是基于令牌桶思想的。并且在此基础上进行了一定程度的扩充，比如`github.com/juju/ratelimit`提供了几种不同特色的令牌桶填充方式：

```go
func NewBucket(fillInterval time.Duration, capacity int64) *Bucket
```

默认的令牌桶，`fillInterval`指每过多长时间向桶里放一个令牌，`capacity`是桶的容量，超过桶容量的部分会被直接丢弃。桶初始是满的。

```go
func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket
```

和普通的`NewBucket()`的区别是，每次向桶中放令牌时，是放`quantum`个令牌，而不是一个令牌。

```go
func NewBucketWithRate(rate float64, capacity int64) *Bucket
```

这个就有点特殊了，会按照提供的比例，每秒钟填充令牌数。例如`capacity`是100，而`rate`是0.1，那么每秒会填充10个令牌。

从桶中获取令牌也提供了几个API：

```go
func (tb *Bucket) Take(count int64) time.Duration {}
func (tb *Bucket) TakeAvailable(count int64) int64 {}
func (tb *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (
    time.Duration, bool,
) {}
func (tb *Bucket) Wait(count int64) {}
func (tb *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {}
```

### 令牌桶原理

从功能上来看，令牌桶模型就是对全局计数的加减法操作过程，可以用buffered channel来完成简单的加令牌取令牌操作。

令牌桶每隔一段固定的时间向桶中放令牌，如果我们记下上一次放令牌的时间为 t1，和当时的令牌数k1，放令牌的时间间隔为ti，每次向令牌桶中放x个令牌，令牌桶容量为cap。现在如果有人来调用`TakeAvailable`来取n个令牌，我们将这个时刻记为t2。在t2时刻，令牌桶中理论上应该有多少令牌呢？伪代码如下：

```go
cur = k1 + ((t2 - t1)/ti) * x
cur = cur > cap ? cap : cur
```

我们用两个时间点的时间差，再结合其它的参数，理论上在取令牌之前就完全可以知道桶里有多少令牌了。那劳心费力地像本小节前面向channel里填充token的操作，理论上是没有必要的。只要在每次`Take`的时候，再对令牌桶中的token数进行简单计算，就可以得到正确的令牌数。是不是很像`惰性求值`的感觉？

在得到正确的令牌数之后，再进行实际的`Take`操作就好，这个`Take`操作只需要对令牌数进行简单的减法即可，记得加锁以保证并发安全。

### 服务瓶颈和QoS

虽然性能指标很重要，但对用户提供服务时还应考虑服务整体的QoS。QoS全称是Quality of Service，顾名思义是服务质量。QoS包含有可用性、吞吐量、时延、时延变化和丢失等指标。一般来讲我们可以通过优化系统，来提高Web服务的CPU利用率，从而提高整个系统的吞吐量。但吞吐量提高的同时，用户体验是有可能变差的。用户角度比较敏感的除了可用性之外，还有时延。虽然你的系统吞吐量高，但半天刷不开页面，想必会造成大量的用户流失。所以在大公司的Web服务性能指标中，除了平均响应时延之外，还会把响应时间的95分位，99分位也拿出来作为性能标准。平均响应在提高CPU利用率没受到太大影响时，可能95分位、99分位的响应时间大幅度攀升了，那么这时候就要考虑提高这些CPU利用率所付出的代价是否值得了。

在线系统的机器一般都会保持CPU有一定的余裕。
