package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/huxiangyu99/advanced-go-programming/ch4rpc/playgrpc/cert/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// server.io *.example.com
	tlsServerName = "server.io"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("../keys/client.pem", "../keys/client.key")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../keys/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   tlsServerName, // NOTE: this is required!
		RootCAs:      certPool,
	})

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
