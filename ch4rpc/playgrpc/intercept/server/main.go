package main

import (
	"context"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/intercept/proto"
	"google.golang.org/grpc"
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

func logFilter(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("logFilter:", info)
	return handler(ctx, req)
}

func panicFilter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("panicFilter:", info)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return handler(ctx, req)
}

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			logFilter, panicFilter,
		)),
	)
	// grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logFilter))
	pb.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
