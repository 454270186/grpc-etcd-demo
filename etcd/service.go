package etcd

import (
	"context"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Register(cli *clientv3.Client, serviceName, addr string) {
	if cli == nil {
		log.Println("etcd client is nil")
		return
	}

	key := fmt.Sprintf("/service/%s", serviceName)
	lease := clientv3.NewLease(cli)

	// Grant a least with TTL
	ctx := context.Background()
	leaseResp, err := lease.Grant(ctx, 10)
	if err != nil {
		log.Fatalln(err)
	}

	// Register the service
	_, err = cli.Put(ctx, key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Fatalln(err)
	}

	// keep the lease alive
	keepAliveChan, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		log.Fatalln(err)
	}

	go func ()  {
		for {
			select {
			case <- ctx.Done():
				return
			case leaseKeepResp := <- keepAliveChan:
				if leaseKeepResp == nil {
					log.Printf("service %s has expried\n", serviceName)
					return
				}
			}
		}	
	}()
}

func Discover(cli *clientv3.Client, serviceName string) (string, error) {
	key := fmt.Sprintf("/service/%s", serviceName)
	ctx := context.Background()

	// Get service addr
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Fatalln(err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service %s not found", serviceName)
	}

	return string(resp.Kvs[0].Value), nil
}