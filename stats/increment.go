package stats

import (
	"go.uber.org/atomic"
	"sync"
	"github.com/funlake/gopkg/timer"
	"github.com/funlake/gopkg/utils/log"
)
type Report struct {
	qps int64
}
type StatsMeta struct {
	ask atomic.Int64
	pass atomic.Int64
	fails atomic.Int64
}
func NewStatsMeta() *StatsMeta{
	return &StatsMeta{}
}
type Increment struct {
	metas sync.Map
}

func (it *Increment) Setup(key string)  {
	if _,set := it.metas.Load(key);!set{
		it.metas.Store(key,NewStatsMeta())
	}
}
func (it *Increment) IncrAsk(key string){
	if meta,set := it.metas.Load(key);set{
		meta.(*StatsMeta).ask.Add(1)
	}
}

type Stats struct {
	prev *StatsMeta
	recent *StatsMeta
	report *Report
	increment *Increment
	ticker *timer.Ticker
	key string
	trigger bool

}
func (s *Stats) rolling(){
	s.ticker.Set(1, s.key, func() {
		if v,ok := s.increment.metas.Load(s.key);ok {
			if !s.trigger {
				s.recent.ask.Store(v.(*StatsMeta).ask.Load())
				s.update(true)
			} else {
				s.prev.ask.Store(v.(*StatsMeta).ask.Load())
				s.update(false)
			}
			s.trigger = !s.trigger
		}
	})
}
func (s *Stats) update(d bool){
	log.Info("recent:%d,prev:%d",s.recent.ask.Load(),s.prev.ask.Load())
	//s.report.ask.Store(s.recent.ask.Load() - s.prev.ask.Load())
	s.report.qps = s.recent.ask.Load() - s.prev.ask.Load()
	if d {
		if s.report.qps < 0 {
			s.report.qps = 0
		}
	} else {
		if s.report.qps < 0 {
			s.report.qps = 0 - s.report.qps
		}
	}
}
