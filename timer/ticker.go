package timer

import "sync"

var cron = NewTimer()
var synconce = sync.Once{}
func NewTicker() *Ticker {
	synconce.Do(func() {
		cron.Ready()
	})
	return &Ticker{slotItems:make(map[string] *SlotItem)}
}
type Ticker struct {
	slotItems map[string] *SlotItem
}
func (this *Ticker) Set(second int,key string,fun func()){
	slotkey := string(second) + "_" + key
	if _,ok := this.slotItems[slotkey];!ok{
		this.slotItems[slotkey] = cron.SetInterval(second,fun)
	}
}
func (this *Ticker) Stop(second int,key string){
	slotkey := string(second) + "_" + key
	if _,ok := this.slotItems[slotkey];ok{
		cron.StopInterval(this.slotItems[slotkey])
		delete(this.slotItems,slotkey)
	}
}
