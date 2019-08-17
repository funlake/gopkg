package cache

import (
	"github.com/pkg/errors"
	"sync"
)

type KvStoreMemory struct {
	localStorage sync.Map
}
func (ksm *KvStoreMemory) Get(key string,opts ... interface{}) (interface{}, error){
	if val,ok := ksm.localStorage.Load(key); ok {
		return val,nil
	}
	return "" ,errors.New("NotFound")
}
func (ksm *KvStoreMemory) Set(key string, val interface{}) (interface{},error) {
	ksm.localStorage.Store(key,val)
	return "" , nil
}
func (ksm *KvStoreMemory) Delete(key string) (interface{},error) {
	ksm.localStorage.Delete(key)
	return "",nil
}
