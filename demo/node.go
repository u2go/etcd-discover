package main

import (
	"etcdDiscover"
	"etcdDiscover/demo/intl"
	"log"
)

func main() {
	discover := etcdDiscover.NewDiscover(intl.Client, intl.ServiceName)
	node, err := etcdDiscover.NewNode("demo")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(discover.Register(node))
}
