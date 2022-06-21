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

	err = client.Call("HelloService.HelloWithoutLogin", "rpc", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)

	err = client.Call("HelloService.Login", "user:password", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)

	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
