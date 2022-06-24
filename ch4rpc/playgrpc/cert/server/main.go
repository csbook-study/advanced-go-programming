package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/cert/proto"
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
	certificate, err := tls.LoadX509KeyPair("../keys/server.pem", "../keys/server.key")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../keys/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		ClientCAs:    certPool,
	})

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
