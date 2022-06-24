package main

import (
	"context"
	"fmt"
	"log"

	"github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/token/auth"
	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/token/proto"
	"google.golang.org/grpc"
)

func main() {
	auth := auth.Authentication{
		User:     "gopher",
		Password: "password",
	}

	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
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
