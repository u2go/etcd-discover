package main

import (
	"etcdDiscover"
	"etcdDiscover/demo/intl"
	"fmt"
	"log"
)

func main() {
	discover := etcdDiscover.NewDiscover(intl.Client, intl.ServiceName)
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
