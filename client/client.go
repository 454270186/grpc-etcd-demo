package main

import (
	"context"
	"log"
	"time"

	"github.com/454270186/grpc-etcd-demo/etcd"
	"github.com/454270186/grpc-etcd-demo/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	//serverAddr = "127.0.0.1:8080"
)

func main() {
	etcdCli := etcd.NewETCDClient()
	defer etcdCli.Close()

	addr, err := etcd.Discover(etcdCli, "hello")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to connect to server", err)
	}
	defer conn.Close()

	// Create a new hello rpc client
	client := hello.NewHelloServiceClient(conn)

	req := hello.HelloReq {
		Name: "yuerfei",
		Age: 20,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	resp, err := client.GetHello(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get RPC Server response: ", resp.GreetMsg)
}

