package main

import (
	"log"
	"net/http"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/pbgo/hello/proto"
)

type HelloService struct{}

func (p *HelloService) Hello(request *pb.String, reply *pb.String) error {
	reply.Value = "hello:" + request.GetValue()
	return nil
}

func main() {
	router := pb.HelloServiceHandler(new(HelloService))
	log.Fatal(http.ListenAndServe(":8080", router))
}
