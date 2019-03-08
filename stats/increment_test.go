package stats

import (
	"testing"
	"time"
	"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/timer"
	"math/rand"
)

func TestStats(t *testing.T) {
	increment := &Increment{}
	increment.Setup("request_baidu")
	go func() {
		stats := &Stats{
			prev  : &StatsMeta{},
			recent: &StatsMeta{},
			report: &Report{},
			increment: increment,
			ticker: timer.NewTicker(),
			key : "request_baidu",
			trigger: false,
		}
		stats.rolling()
		tm := timer.NewTimer()
		tm.Ready()
		tm.SetInterval(1, func() {
			log.Info("qps : %d",stats.report.qps)
		})
	}()

	for  {
		increment.IncrAsk("request_baidu")
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Millisecond * time.Duration(r.Intn(100)))
	}

}
