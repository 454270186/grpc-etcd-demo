package etcd

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)


const (
	etcdAddr = "localhost:2379"
)

func NewETCDClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln("Fail to connect to etcd")
	}

	return cli
}