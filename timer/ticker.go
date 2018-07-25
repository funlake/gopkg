package timer
var cron = NewTimer()
func NewTicker() *Ticker {
	cron.Ready()
	return &Ticker{slotItems:make(map[string] *SlotItem)}
}
type Ticker struct {
	slotItems map[string] *SlotItem
}
func (this *Ticker) Set(second int,key string,fun func()){
	if _,ok := this.slotItems[key];!ok{
		this.slotItems[key] = cron.SetInterval(second,fun)
	}
}
func (this *Ticker) Stop(key string){
	if _,ok := this.slotItems[key];ok{
		cron.StopInterval(this.slotItems[key])
		delete(this.slotItems,key)
	}
}
