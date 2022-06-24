package main

import (
	"context"
	"log"
	"net"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/publickey/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HelloServiceImpl struct {
	pb.UnimplementedHelloServiceServer
}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *pb.String,
) (*pb.String, error) {
	reply := &pb.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {
	creds, err := credentials.NewServerTLSFromFile("../keys/server.crt", "../keys/server.key")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
