package etcdDiscover

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	client "go.etcd.io/etcd/client/v3"
	"strconv"
)

var (
	// etcd 的 key 前缀
	etcdKeyPrefix = "ed/"
	// NodeTtl 节点超时时间，单位秒
	NodeTtl int64 = 3
)

// NewDiscover 创建服务发现实例
func NewDiscover(client *client.Client, serviceName string) *discover {
	return &discover{
		client: client,

		serviceName: serviceName,
	}
}

type discover struct {
	client *client.Client

	// 服务名
	serviceName string
}

// Register 节点注册
// 出错返回 error，返回error的时候要重启服务
func (d *discover) Register(node *Node) error {
	var key = fmt.Sprintf("%s%s/%d", etcdKeyPrefix, d.serviceName, node.Id)
	b, err := json.Marshal(node)
	if err != nil {
		return err
	}
	var value = string(b)

	lease := client.NewLease(d.client)
	leaseRes, err := lease.Grant(Ctx(), NodeTtl)
	if err != nil {
		return err
	}
	leaseId := leaseRes.ID
	defer lease.Revoke(Ctx(), leaseId)

	// 设置 key
	kv := client.NewKV(d.client)
	_, err = kv.Put(Ctx(), key, value, client.WithLease(leaseId))
	if err != nil {
		return err
	}

	keepRespChan, err := lease.KeepAlive(Ctx(), leaseId)
	if err != nil {
		return err
	}

	// 续约应答
	for {
		select {
		case <-keepRespChan:
			if keepRespChan == nil {
				goto END
			}
		}
	}
END:

	return errors.New("lease closed")
}

// WatchNodes 监听节点变化。
// handle会在收到节点变化时被调用，会调用多次
func (d *discover) WatchNodes(handle func(nodes map[int64]*NodeAction)) error {
	key := fmt.Sprintf("%s%s/", etcdKeyPrefix, d.serviceName)
	watcher := client.NewWatcher(d.client)
	watchChan := watcher.Watch(Ctx(), key, client.WithPrefix())
	for res := range watchChan {
		nodes := map[int64]*NodeAction{}
		for _, event := range res.Events {
			node := &Node{}
			if event.Kv.Value == nil {
				idStr := string(event.Kv.Key[bytes.LastIndexByte(event.Kv.Key, '/')+1:])
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					return err
				}
				node.Id = id
			} else {
				err := json.Unmarshal(event.Kv.Value, node)
				if err != nil {
					return err
				}
			}
			nodes[node.Id] = &NodeAction{
				Action: event.Type,
				Node:   node,
			}
		}
		handle(nodes)
	}
	return nil
}

// GetNodes 获取所有活跃节点
func (d *discover) GetNodes() (map[int64]*Node, error) {
	key := fmt.Sprintf("%s%s/", etcdKeyPrefix, d.serviceName)
	kv := client.NewKV(d.client)
	res, err := kv.Get(Ctx(), key, client.WithPrefix())
	if err != nil {
		return nil, err
	}
	if res.Count == 0 {
		return nil, nil
	}
	var nodes = map[int64]*Node{}
	for _, v := range res.Kvs {
		var node = &Node{}
		err := json.Unmarshal(v.Value, node)
		if err != nil {
			return nil, err
		}
		nodes[node.Id] = node
	}
	return nodes, nil
}
