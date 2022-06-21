package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	rpc.RegisterName("KVStoreService", NewKVStoreService())

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
