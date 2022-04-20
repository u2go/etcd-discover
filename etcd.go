package etcdDiscover

import (
	"encoding/json"
	client "go.etcd.io/etcd/client/v3"
	"time"
)

// NewEtcd 创建服务发现实例
func NewEtcd(etcdEndpoints []string) (*Etcd, error) {
	client0, err := client.New(client.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Etcd{
		Client: client0,
	}, nil
}

type Etcd struct {
	Client *client.Client
}

// Watch 监听节点变化。
// handle 会在收到节点变化时被调用，会调用多次
func (d *Etcd) Watch(key string, isDir bool, handle func(res client.WatchResponse)) {
	watcher := client.NewWatcher(d.Client)
	var opts []client.OpOption
	if isDir {
		opts = append(opts, client.WithPrefix())
	}
	watchChan := watcher.Watch(Ctx(), key, opts...)
	for res := range watchChan {
		handle(res)
	}
}

// Get 获取Key前缀的所有节点
func (d *Etcd) Get(key string, isDir bool) (res *client.GetResponse, err error) {
	kv := client.NewKV(d.Client)
	var opts []client.OpOption
	if isDir {
		opts = append(opts, client.WithPrefix())
	}
	return kv.Get(Ctx(), key, opts...)
}

// Put 设置内容
func (d *Etcd) Put(key string, valueRaw any) (err error) {
	value, err := json.Marshal(valueRaw)
	if err != nil {
		return err
	}
	kv := client.NewKV(d.Client)
	_, err = kv.Put(Ctx(), key, string(value))
	return err
}
