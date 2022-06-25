package main

import (
	"context"
	"net"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/grpcprotobuf/rest/proto"
	"google.golang.org/grpc"
)

type RestServiceImpl struct {
	pb.UnimplementedRestServiceServer
}

func (r *RestServiceImpl) Get(ctx context.Context, message *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{Value: "Get hi:" + message.Value + "#"}, nil
}

func (r *RestServiceImpl) Post(ctx context.Context, message *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{Value: "Post hi:" + message.Value + "@"}, nil
}
func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterRestServiceServer(grpcServer, new(RestServiceImpl))
	lis, _ := net.Listen("tcp", ":5000")
	grpcServer.Serve(lis)
}
