package main

import (
	"fmt"
	etcdDiscover "github.com/u2go/etcd-discover"
	"log"
)

func main() {
	discover, err := etcdDiscover.NewDiscover([]string{"127.0.0.1:2379"}, "demo")
	if err != nil {
		log.Fatalln(err)
	}
	// 所有节点
	log.Println(discover.GetNodes())
	// 监听节点变化
	log.Fatalln(discover.WatchNodes(func(nodes map[int64]*etcdDiscover.NodeAction) {
		fmt.Println("node change")
		for _, node := range nodes {
			fmt.Println(node.Action, node.Node.Id, node.Node.Name, node.Node.HostName, node.Node.Addr)
		}
		// 当前所有节点
		nodes1, err := discover.GetNodes()
		fmt.Println("all nodes", nodes1, err)
	}))
}
