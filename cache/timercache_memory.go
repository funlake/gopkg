package cache

import "github.com/pkg/errors"

type TimerCacheMemory struct {
	store *KvStoreMemory
}

func (tcm *TimerCacheMemory) Flush(key string){
	_, _ = tcm.store.Delete(key)
}
func (tcm *TimerCacheMemory) Get(hk string, k string, wheel int) (string, error){
	if v,ok := tcm.store.localStorage.Load(hk+k);ok{
		return v.(string),nil
	}
	return "",errors.New("NotFound")
}
func (tcm *TimerCacheMemory) SetStore(store KvStore){
	tcm.store = store.(*KvStoreMemory)
}
func (tcm *TimerCacheMemory) GetStore() KvStore{
	return tcm.store
}