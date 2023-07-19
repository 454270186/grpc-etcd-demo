package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/454270186/grpc-etcd-demo/etcd"
	"github.com/454270186/grpc-etcd-demo/hello"
	"google.golang.org/grpc"
)

const (
	addr = "127.0.0.1:8080"
)

type helloServiceServer struct {
	hello.UnimplementedHelloServiceServer
}

func (h *helloServiceServer) GetHello(ctx context.Context, req *hello.HelloReq) (*hello.HelloRes, error) {
	greet := fmt.Sprintf("hello %s, age: %d", req.Name, req.Age)
	return &hello.HelloRes{
		Code: 0,
		GreetMsg: greet,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	
	server := grpc.NewServer()
	hello.RegisterHelloServiceServer(server, &helloServiceServer{})

	// Register service
	etcdCli := etcd.NewETCDClient()
	defer etcdCli.Close()
	etcd.Register(etcdCli, "hello", addr)

	log.Println("RPC Server starting listening at", addr)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}