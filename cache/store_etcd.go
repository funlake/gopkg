package cache

import (
  "context"
  "go.etcd.io/etcd/clientv3"
  "go.etcd.io/etcd/clientv3/concurrency"
  "strings"
  "time"
)

func NewKvStoreEtcd() *KvStoreEtcd {
  return &KvStoreEtcd{}
}

type KvStoreEtcd struct {
   conn *clientv3.Client
}

func (es *KvStoreEtcd) Connect(dns,pwd string){
  var err error
  es.conn,err = clientv3.New(clientv3.Config{
    Endpoints:strings.Split(dns,","),
    DialTimeout: time.Second * 3,
  })
  if err != nil {
    panic("No available etcd server")
  }
}
func (es *KvStoreEtcd) Get(key string) (interface{},error){
  return es.conn.Get(context.TODO(),key)
}
func (es *KvStoreEtcd) Set(key string , val interface{}){
  concurrency.NewSTM(es.conn, func(stm concurrency.STM) error {
    stm.Put(key,val.(string))
    return nil
  })
}
func (es *KvStoreEtcd) 	GetPool() interface{}{
  return es.conn
}
func (es *KvStoreEtcd) Delete(key string) {
  es.conn.Delete(context.TODO(),key)
}
func (es *KvStoreEtcd) Watch(ctx context.Context,key string) (clientv3.WatchChan) {
  return es.conn.Watch(ctx,key)
}
