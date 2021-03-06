package cache

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	// cv3  "go.etcd.io/etcd/clientv3"
	cv3 "github.com/coreos/etcd/clientv3"
	"strings"
	"time"
)

func NewKvStoreEtcd() *KvStoreEtcd {
	return &KvStoreEtcd{}
}

type KvStoreEtcd struct {
	conn *cv3.Client
}

func (es *KvStoreEtcd) Connect(dsn, pwd string) {
	var err error
	es.conn, err = cv3.New(cv3.Config{
		Endpoints:   strings.Split(dsn, ","),
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		panic("No available etcd server:" + err.Error())
	}
}
func (es *KvStoreEtcd) ConnectWithTls(dsn, tlsc interface{}) error {
	var err error
	es.conn, err = cv3.New(cv3.Config{
		Endpoints:   strings.Split(dsn.(string), ","),
		DialTimeout: time.Second * 3,
		TLS:         tlsc.(*tls.Config),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	return err
}

//todo : cancel context needed,cv3.OpOption
func (es *KvStoreEtcd) Get(key string, opts ...interface{}) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	var opOptions []cv3.OpOption
	j:=len(opts)
	for i:=0;i<j;i++{
		opOptions = append(opOptions, opts[i].(cv3.OpOption))
	}
	return es.conn.Get(ctx, key, opOptions...)
}
func (es *KvStoreEtcd) Set(key string, val interface{}) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	return es.conn.Put(ctx, key, val.(string))
}
func (es *KvStoreEtcd) HashSet(hk, key string, val interface{}) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	return es.conn.Put(ctx, hk+"/"+key, val.(string))
}
func (es *KvStoreEtcd) GetPool() interface{} {
	return es.conn
}
func (es *KvStoreEtcd) Delete(key string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	return es.conn.Delete(ctx, key)
}
func (es *KvStoreEtcd) Watch(ctx context.Context, key string, opts ...cv3.OpOption) cv3.WatchChan {
	return es.conn.Watch(ctx, key, opts...)
}
func (es *KvStoreEtcd) GetActiveCount() int {
	return 1
	//new version will have it?
	//return int(es.conn.ActiveConnection().ChannelzMetric().CallsStarted)
}
