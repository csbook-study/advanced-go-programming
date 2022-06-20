package main

import (
	"log"
	"net"
	"net/rpc"
)

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

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("listener.Accept error:", err)
	}

	rpc.ServeConn(conn)
}
