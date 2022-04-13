package main

import (
	etcdDiscover "github.com/u2go/etcd-discover"
	"log"
)

func main() {
	discover, err := etcdDiscover.NewDiscover([]string{"127.0.0.1:2379"}, "demo")
	if err != nil {
		log.Fatalln(err)
	}
	node, err := etcdDiscover.NewNode("demo-node-1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(discover.Register(node))
}
