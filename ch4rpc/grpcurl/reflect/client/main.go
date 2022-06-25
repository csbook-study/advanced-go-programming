package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/introduce2grpc/helloworld/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &pb.String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}
