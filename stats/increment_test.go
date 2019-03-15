package stats

import (
	"testing"
	"time"
	"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/timer"
)

func TestStats(t *testing.T) {
	increment := &Increment{}
	increment.Setup("request_baidu")
	go func() {
		tm := timer.NewTimer()
		tm.Ready()
		tm.SetInterval(1, func() {
			log.Info("qps : %d",increment.GetStat("request_baidu").GetReport().GetQps())
		})
	}()

	for  {
		increment.IncrRequest("request_baidu")
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Millisecond * 300)
	}

}
