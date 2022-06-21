package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/introduce2grpc/pubsub/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewPubsubServiceClient(conn)

	reply, err := client.Publish(
		context.Background(), &pb.String{Value: "golang: hello Go"},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("publish golang", reply.GetValue())
	reply, err = client.Publish(
		context.Background(), &pb.String{Value: "docker: hello Docker"},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("publish docker", reply.GetValue())
}
