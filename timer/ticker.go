package timer

import "sync"

var (
	cron   *timer
	once   = sync.Once{}
	ticker *Ticker
)

func NewTicker() *Ticker {
	once.Do(func() {
		cron = NewTimer()
		cron.Ready()
		ticker = &Ticker{}
	})
	return ticker
}

type Ticker struct {
	slotItems sync.Map
}

func (t *Ticker) Set(second int, key string, fun func()) {
	slotkey := string(second) + "_" + key
	if _, ok := t.slotItems.Load(slotkey); !ok {
		t.slotItems.Store(slotkey,cron.SetInterval(second, fun))
	}
}
func (t *Ticker) Stop(second int, key string) {
	slotkey := string(second) + "_" + key
	if tk, ok := t.slotItems.Load(slotkey); ok {
		cron.StopInterval(tk.(*SlotItem))
		t.slotItems.Delete(slotkey)
	}
}
