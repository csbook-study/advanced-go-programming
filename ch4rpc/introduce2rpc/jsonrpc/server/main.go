package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 模拟客户端请求
// echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
// {"id":1,"result":"hello, hello","error":null}

type HelloService struct{}

func (p *HelloService) Hello(request string, response *string) error {
	*response = "hello, " + request
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("net.Listen error: ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
