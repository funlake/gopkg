package breaker

import (
	"testing"

	"net/http"
	"github.com/funlake/gopkg/utils/log"
	"time"
)

func TestBreaker_Run(t *testing.T) {
	bre := NewBreaker("hello",2,3,30)
	for {
		bre.Run(func() {
			res,_ := http.Get("http://admin.etcchebao.com/index2.php")
			log.Success("%d",res.StatusCode)
			//now := time.Now()
			//rand.Seed(time.Now().UnixNano())
			//t := rand.Intn(5)
			//time.Sleep(time.Second * time.Duration(t))
			//log.Info("%s",time.Since(now))
		}, func() {
			log.Success("goog job")
		}, func() {
			log.Error("break!")
		})
		time.Sleep(time.Millisecond * 100)
	}
}
