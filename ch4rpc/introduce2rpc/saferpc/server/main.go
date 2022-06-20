package main

import (
	"log"
	"net"
	"net/rpc"
)

const HelloServiceName = "HelloService"

type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

type HelloService struct{}

func (p *HelloService) Hello(request string, response *string) error {
	*response = "hello, " + request
	return nil
}

func main() {
	RegisterHelloService(new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("net.Listen error: ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("listener.Accept error: ", err)
		}

		go rpc.ServeConn(conn)
	}
}
