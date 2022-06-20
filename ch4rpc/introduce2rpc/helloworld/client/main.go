package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("rpc.Dial error: ", err)
	}

	var reply string
	err = client.Call("HelloService.Hello", "rpc", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
