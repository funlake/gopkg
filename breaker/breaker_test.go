package breaker

import (
	"testing"
	//. "github.com/smartystreets/goconvey/convey"
	//"net/http"
	"github.com/funlake/gopkg/utils/log"
	"math/rand"
	"time"
)

//func TestBreaker_Run(t *testing.T) {
//	bre := NewBreaker("hello",2,30,10)
//	for {
//		bre.Run(func() {
//			res,_ := http.Get("http://www.google.com")
//			log.Success("%d",res.StatusCode)
//			//now := time.Now()
//			//rand.Seed(time.Now().UnixNano())
//			//t := rand.Intn(5)
//			//time.Sleep(time.Second * time.Duration(t))
//			//log.Info("%s",time.Since(now))
//		}, func() {
//			log.Success("goog job")
//		}, func(run bool) {
//			log.Error("break!")
//		})
//		time.Sleep(time.Millisecond * 100)
//	}
//}
func BenchmarkBreaker_Run(b *testing.B) {
	b.SetParallelism(3)
	bre := NewBreaker("simulate break request")
	bre.SetTimeout(2)
	bre.SetRate(50)
	bre.SetWindow(3)
	bre.SetMin(3)
	b.N = 1000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bre.Run(func() {
				//http.Get("http://www.google.com")
				//time.Sleep(3 * time.Second)
				//log.Info("what the hell?")
				//log.Success("%d",res.StatusCode)
				//now := time.Now()
				rand.Seed(time.Now().UnixNano())
				t := rand.Intn(5)
				time.Sleep(time.Second * time.Duration(t))
				//log.Info("%s",time.Since(now))
			}, func() {
				log.Success("job pass")
			}, func(run bool) {
				log.Error("break!")
			})
			//log.Success("what ?")
		}
	})
	//time.Sleep(time.Second * 20)
}

//func TestNewBreaker(t *testing.T) {
//	bre := NewBreaker("hello")
//	bre.SetTimemout(2)
//	bre.SetRate(10)
//	bre.SetWindow(30)
//	bre.SetMin(3)
//	Convey("new breaker",t, func() {
//		So(bre,ShouldNotBeNil)
//	})
//	Convey("twice new calling",t,func(){
//		c := NewBreaker("hello2")
//		bre.SetTimemout(2)
//		bre.SetRate(10)
//		bre.SetWindow(30)
//		bre.SetMin(3)
//		So(bre,ShouldEqual,c)
//	})
//}
//
//func TestBreakerPassBreaker(t *testing.T) {
//	b := NewBreaker("hello",2,30,10)
//	for i:=1;i<=100;i++{
//		go func(i int) {
//			b.pass = b.pass + i
//		}(i)
//	}
//	time.Sleep(time.Second * 1)
//	Convey("roundtine object counter",t, func() {
//		So(b.pass,ShouldEqual,5050)
//	})
//}
