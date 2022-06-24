package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/publickey/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(
		"../keys/server.crt", "server.grpc.io",
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost:1234", grpc.WithTransportCredentials(creds))
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
