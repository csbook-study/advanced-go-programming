package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 模拟服务端
// nc -l 1234
// {"method":"HelloService.Hello","params":["hello"],"id":0}

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("net.Dial error: ", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	err = client.Call("HelloService.Hello", "json rpc", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
