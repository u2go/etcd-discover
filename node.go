package etcdDiscover

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	"time"
)

const (
	NodeActionPut    = mvccpb.PUT
	NodeActionDelete = mvccpb.DELETE
)

func NewNode(name string) (*Node, error) {
	ip, err := Ip()
	if err != nil {
		return nil, err
	}
	hostname, err := Hostname()
	if err != nil {
		return nil, err
	}
	return &Node{
		Id:       time.Now().UnixNano(),
		Name:     name,
		Addr:     ip,
		HostName: hostname,
	}, nil
}

type Node struct {
	Id       int64
	Name     string
	Addr     string
	HostName string
}

type NodeAction struct {
	Action mvccpb.Event_EventType
	Node   *Node
}
