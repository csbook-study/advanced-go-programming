package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/moby/moby/pkg/pubsub"
	"google.golang.org/grpc"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/introduce2grpc/pubsub/proto"
)

type PubsubServiceImpl struct {
	pub *pubsub.Publisher

	pb.UnimplementedPubsubServiceServer
}

func NewPubsubServiceImpl() *PubsubServiceImpl {
	return &PubsubServiceImpl{
		pub: pubsub.NewPublisher(100*time.Millisecond, 10),
	}
}

func (p *PubsubServiceImpl) Publish(
	ctx context.Context, arg *pb.String,
) (*pb.String, error) {
	p.pub.Publish(arg.GetValue())
	return &pb.String{}, nil
}

func (p *PubsubServiceImpl) Subscribe(
	arg *pb.String, stream pb.PubsubService_SubscribeServer,
) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := stream.Send(&pb.String{Value: v.(string)}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterPubsubServiceServer(grpcServer, NewPubsubServiceImpl())

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
