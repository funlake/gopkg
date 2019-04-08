package cache

import (
  "context"
  "errors"
  "fmt"
  "github.com/funlake/gopkg/utils"
  "go.etcd.io/etcd/clientv3"
  "sync"
)

func NewTimerCacheEtcd() *TimerCacheEtcd {
  return &TimerCacheEtcd{}
}
type TimerCacheEtcd struct {
  store *KvStoreEtcd
  local sync.Map
}
func (tc *TimerCacheEtcd) Get(hk string,k string,wheel int) (string,error){
  var rv string
  val,ok := tc.local.Load(k)
  if !ok {
    resp, err := tc.store.Get(k)
    if err == nil {
      for _, e := range resp.(*clientv3.GetResponse).Kvs {
        rv = string(e.Value)
      }
    } else {
      return "",err
    }
    if rv != "" {
      utils.WrapGo(func() {
        tc.Watch(k)
      },fmt.Sprintf("watch-key-%s",k))
      tc.local.Store(k, rv)
      return rv,nil
    }else{
      tc.local.Store(k, "")
      return "",errors.New("Value Not set")
    }
  }
  rv = val.(string)
  return rv,nil
}
func (tc *TimerCacheEtcd) Flush(){
  tc.local.Range(func(key, value interface{}) bool {
    tc.local.Delete(key)
    return true
  })
}
func (tc *TimerCacheEtcd) SetStore(store *KvStoreEtcd){
  tc.store = store
}
func (tc *TimerCacheEtcd) Watch(key string){
  ctx,cancel := context.WithCancel(context.Background())
  wc := tc.store.Watch(ctx,key)
  for v := range wc {
    if v.Err() != nil{
      panic(v.Err().Error())
    }
    for _,e := range v.Events{
      tp := fmt.Sprintf("%v",e.Type)
      switch tp {
        case "DELETE" :
          tc.local.Store(key,"")
          cancel()
        break
        case "PUT" :
          tc.local.Store(key,string(e.Kv.Value))
        break
      }
    }
  }
}
