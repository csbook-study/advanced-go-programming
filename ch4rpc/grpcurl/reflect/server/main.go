package main

import (
	"context"
	"log"
	"net"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/grpcurl/reflect/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
