package cache

import (
  "context"
  "crypto/tls"
  "github.com/funlake/gopkg/utils/log"
  cv3  "go.etcd.io/etcd/clientv3"
  con "go.etcd.io/etcd/clientv3/concurrency"
  "strings"
  "time"
)

func NewKvStoreEtcd() *KvStoreEtcd {
  return &KvStoreEtcd{}
}

type KvStoreEtcd struct {
   conn *cv3.Client
}

func (es *KvStoreEtcd) Connect(dsn,pwd string){
  var err error
  es.conn,err = cv3.New(cv3.Config{
    Endpoints:strings.Split(dsn,","),
    DialTimeout: time.Second * 3,
  })
  if err != nil {
    panic("No available etcd server:"+err.Error())
  }
}
func (es *KvStoreEtcd) ConnectWithTls(dsn ,tlsc interface{})(error){
  var err error
  es.conn,err = cv3.New(cv3.Config{
    Endpoints:strings.Split(dsn.(string),","),
    DialTimeout: time.Second * 3,
    TLS: tlsc.(*tls.Config),
  })
  return err
}
//todo : cancel context needed
func (es *KvStoreEtcd) Get(key string) (interface{},error){
  ctx,_ := context.WithTimeout(context.Background(),time.Millisecond * 500)
  r,err := es.conn.Get(ctx,key)
  return r,err
}
func (es *KvStoreEtcd) Set(key string , val interface{}){
  _,err := con.NewSTM(es.conn, func(stm con.STM) error {
    stm.Put(key,val.(string))
    return nil
  })
  if err != nil{
    log.Error(err.Error())
  }
}
func (es *KvStoreEtcd) HashSet(hk,key string , val interface{})(string,error){
  _,err := con.NewSTM(es.conn, func(stm con.STM) error {
    stm.Put(hk+"/"+key,val.(string))
    return nil
  })
  if err != nil{
    log.Error(err.Error())
  }
  return "",err
}
func (es *KvStoreEtcd) 	GetPool() interface{}{
  return es.conn
}
func (es *KvStoreEtcd) Delete(key string) {
  _,err := es.conn.Delete(context.TODO(),key)
  if err != nil {
    log.Error(err.Error())
  }
}
func (es *KvStoreEtcd) Watch(ctx context.Context,key string) (cv3.WatchChan) {
  return es.conn.Watch(ctx,key)
}
func (es *KvStoreEtcd) GetActiveCount() int{
  return int(es.conn.ActiveConnection().ChannelzMetric().CallsStarted)
}
