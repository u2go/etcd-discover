package intl

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var (
	Client      *clientv3.Client
	ServiceName = "demo"
)

func init() {
	// 建立连接
	var err error
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
