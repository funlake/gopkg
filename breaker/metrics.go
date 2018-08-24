package breaker

import "sync"
var (
	metrcSync sync.Once
	metrices *Metrics
)
func NewMetrics() *Metrics{
	metrcSync.Do(func() {
		metrices = &Metrics{}
	})
	return metrices
}
type Metrics struct{
	entities sync.Map
}
func (m *Metrics) NewEntity(id string,timeout int,window int) *MestricsEntity{
	me := m.GetEntity(id)
	if me == nil{
		me = &MestricsEntity{timeout:timeout,window:window,pass:0,broken:0}
		m.entities.Store(id,me)
		return me
	}
	return me
}
func (m *Metrics)GetEntity(id string) *MestricsEntity{
	if e,ok := m.entities.Load(id);ok{
		return e.(*MestricsEntity)
	}
	return nil

}
type MestricsEntity struct{
	timeout int
	window int
	pass int
	broken int
}

func (me *MestricsEntity) GetTimeout() int{
	return me.timeout
}

func (me *MestricsEntity) GetWindow() int{
	return me.window
}
