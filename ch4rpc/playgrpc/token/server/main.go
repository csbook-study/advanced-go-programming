package main

import (
	"context"
	"log"
	"net"

	"github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/token/auth"
	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/token/proto"
	"google.golang.org/grpc"
)

type HelloServiceImpl struct {
	auth *auth.Authentication

	pb.UnimplementedHelloServiceServer
}

func NewHelloServiceImpl() *HelloServiceImpl {
	return &HelloServiceImpl{
		auth: &auth.Authentication{
			User:     "gopher",
			Password: "password",
		},
	}
}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *pb.String,
) (*pb.String, error) {
	if err := p.auth.Auth(ctx); err != nil {
		return nil, err
	}

	reply := &pb.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, NewHelloServiceImpl())

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
