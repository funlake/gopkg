package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var timerEtcdCache *TimerCacheEtcd
var cacheSyncOnce2 = &sync.Once{}

func initEtcdTest() {
	cacheSyncOnce2.Do(func() {
		timerEtcdCache = NewTimerCacheEtcd()
		etcdStore := NewKvStoreEtcd()
		etcdStore.Connect("127.0.0.1:2379", "")
		timerEtcdCache.SetStore(etcdStore)
	})
}

func TestTimerCacheEtcd_Get(t *testing.T) {
	initEtcdTest()
	for i := 0; i < 10000; i++ {
		val, err := timerEtcdCache.Get("", "hello/world", 0)
		if err != nil {
			t.Log(err.Error())
			fmt.Printf("%s\n", err.Error())
		} else {
			fmt.Printf("%s\n", val)
		}
		time.Sleep(time.Second * 1)
	}
}
