package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tmc = &TimerCacheMemory{}
func TestTimerCacheMemory_SetStore(t *testing.T) {
	tmc.SetStore(&KvStoreMemory{})
}
func TestTimerCacheMemory_GetStore(t *testing.T) {
	assert.NotNil(t,tmc.GetStore())
}
func TestTimerCacheMemory_Get(t *testing.T) {
	_, _ = tmc.GetStore().(*KvStoreMemory).Set("k1", "hello")
	v, _ := tmc.Get("", "k1", 0)
	assert.Equal(t,"hello",v)
}

func TestTimerCacheMemory_Flush(t *testing.T) {
	_, _ = tmc.GetStore().(*KvStoreMemory).Set("k2", "world")
	v, _:= tmc.Get("","k2",0)
	assert.Equal(t,"world",v)
	tmc.Flush("k2")
	_, err := tmc.Get("","k2",0)
	assert.Error(t,err)
}