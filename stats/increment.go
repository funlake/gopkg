package stats

import (
	"go.uber.org/atomic"
	"sync"
	"github.com/funlake/gopkg/timer"
)
type Report struct {
	qps int64
}
type StatsMeta struct {
	request atomic.Int64
	pass    atomic.Int64
	fails   atomic.Int64
}
func NewStatsMeta() *StatsMeta{
	return &StatsMeta{}
}
type Increment struct {
	metas sync.Map
	stats sync.Map
}

func (it *Increment) Setup(key string) bool {
	if _,set := it.metas.Load(key);!set{
		it.metas.Store(key,NewStatsMeta())
		it.initStat(key)
		return false
	}
	return true
}
func (it *Increment) GetStat(key string) *Stats{
	if stat,set := it.stats.Load(key);set{
		return stat.(*Stats)
	}
	return nil
}
func (it *Increment) GetAllStats() sync.Map{
	return it.stats
}
func (it *Increment) initStat(key string){
	stats := &Stats{
		prev  : &StatsMeta{},
		recent: &StatsMeta{},
		report: &Report{},
		increment: it,
		ticker: timer.NewTicker(),
		key : key,
		trigger: false,
	}
	stats.rolling()
	it.stats.Store(key,stats)
}
func (it *Increment) IncrRequest(key string){
	if meta,set := it.metas.Load(key);set{
		meta.(*StatsMeta).request.Add(1)
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
				s.recent.request.Store(v.(*StatsMeta).request.Load())
				s.update(true)
			} else {
				s.prev.request.Store(v.(*StatsMeta).request.Load())
				s.update(false)
			}
			s.trigger = !s.trigger
		}
	})
}
func (s *Stats) update(d bool){
	s.report.qps = s.recent.request.Load() - s.prev.request.Load()
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
func (s *Stats) GetReport() *Report{
	return s.report
}
func (r *Report) GetQps() int64{
	return r.qps
}
