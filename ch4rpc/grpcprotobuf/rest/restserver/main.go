package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/grpcprotobuf/rest/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := pb.RegisterRestServiceHandlerFromEndpoint(
		ctx, mux, "localhost:5000",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", mux)
}
