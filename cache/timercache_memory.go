package cache

type TimerCacheMemory struct {
	store *KvStoreMemory
}

func (tcm *TimerCacheMemory) Flush(key string){

}
func (tcm *TimerCacheMemory) Get(hk string, k string, wheel int) (string, error){
	return "",nil
}
func (tcm *TimerCacheMemory) SetStore(store KvStore){
	tcm.store = store.(*KvStoreMemory)
}
func (tcm *TimerCacheMemory) GetStore() KvStore{
	return tcm.store
}