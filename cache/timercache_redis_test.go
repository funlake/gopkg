package cache

import (
	"testing"
	"os"
)
var timercache *TimerCacheRedis
func initTest(){
	timercache = NewTimerCacheRedis()
	redisStore := NewKvStoreRedis()
	redisStore.Connect(os.Getenv("REDIS_ETC1_HOST")+":"+os.Getenv("REDIS_ETC1_PORT"),os.Getenv("REDIS_ETC1_PASSWORD"))
	//log.Error("%s",redisStore.pool.Get().Do("get","abc"))
	timercache.SetStore(redisStore)
}
func BenchmarkTimerCacheRedis_Get(b *testing.B) {
	initTest()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			//hashget,如果是普通get，将第一个参数置空即可
			_,err := timercache.Get("gateway:proxy","access_token",3)
			if err != nil{
				b.Error(err.Error())
				break
			}
		}
	})
}

func TestTimerCacheRedis_Get(t *testing.T) {
	initTest()
	_,err := timercache.Get("gateway:proxy","access_token",3)
	if err != nil{
		t.Error(err.Error())
	}
}
