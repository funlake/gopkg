package cache
//本地介于redis之上加一层缓存
//自动定时更新缓存
//todo : lru 实现
import (
	"sync"
	"github.com/funlake/gopkg/timer"
	"github.com/funlake/gopkg/utils/log"
	"strings"
)
type TimerCacheRedis struct{
	mu sync.Mutex
	store *KvStoreRedis
	local map[string] string
	ticker *timer.Ticker
}

func NewTimerCacheRedis() *TimerCacheRedis{
	return &TimerCacheRedis{local:make(map[string] string),ticker:timer.NewTicker()}
}
func (tc *TimerCacheRedis) Flush(){
	tc.mu.Lock()
	defer tc.mu.Unlock()
	for k := range tc.local{
		delete(tc.local,k)
		//ticker.Stop(k)
	}
}
func (tc *TimerCacheRedis) SetStore(store *KvStoreRedis){
	tc.store = store
}
func (tc *TimerCacheRedis) GetStore() *KvStoreRedis{
	return tc.store
}
func (tc *TimerCacheRedis) Get(hk string,k string,wheel int) (string,error){
	tc.mu.Lock()
	defer tc.mu.Unlock()
	localCacheKey := hk+"_"+k
	if _,ok := tc.local[localCacheKey];ok{
		return tc.local[localCacheKey],nil
	}else{
		log.Info("Access redis for setting : %s_%s",hk,k)
		v, err := tc.store.HashGet(hk, k)
		if err == nil || strings.Contains(err.Error(),"nil returned") {
			//tc.local[localCacheKey] = v.(string)
			log.Info("set cache value even cache is empty :%s",localCacheKey)
			tc.local[localCacheKey] = v.(string)
			tc.ticker.Set(wheel,localCacheKey, func() {
				tc.mu.Lock()
				defer tc.mu.Unlock()
				//log.Info("每%d秒定时检查%s",wheel,localCacheKey)
				v,err := tc.store.HashGet(hk,k)
				//假如redis服务器挂了,得保留之前的本地缓存值
				if err != nil {
					if strings.Contains(err.Error(),"nil returned"){
						log.Error("Empty value deteced(%d s) : %s,remove ticker run: error:%s",wheel,localCacheKey,err.Error())
						//如果值返回为空，则停掉定时更新器
						tc.ticker.Stop(wheel,localCacheKey)
						//赋空值,如要情况缓存，可调用/api-cleancache接口
						tc.local[localCacheKey] = v.(string)
						//delete(tc.local,localCacheKey)
					}else{
						//发生redis连接故障，则继续保持旧有缓存
						log.Error("Redis seems has gone,we do not clear cache if redis is down to keep gateway service's running")
					}
					//delete(tc.local, localCacheKey)
				}else{
					//log.Warning("更新数据%s : %s",localCacheKey,v.(string))
					tc.local[localCacheKey] = v.(string)
				}
			})
			return tc.local[localCacheKey],nil
		}else{
			log.Error(err.Error())
			return "",err
		}
	}
	//log.Warning("waht the fuck?")
	return tc.local[localCacheKey],nil
}