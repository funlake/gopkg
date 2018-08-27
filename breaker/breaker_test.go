package breaker

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
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
	b.SetParallelism(20)
	bre := NewBreaker("request google",2,30,10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			bre.Run(func() {
				 http.Get("http://www.google.com")
				//log.Success("%d",res.StatusCode)
				//now := time.Now()
				//rand.Seed(time.Now().UnixNano())
				//t := rand.Intn(5)
				//time.Sleep(time.Second * time.Duration(t))
				//log.Info("%s",time.Since(now))
			}, func() {
				//log.Success("goog job")
			}, func(run bool) {
				//log.Error("break!")
			})
		}
	})

}

func TestNewBreaker(t *testing.T) {
	b := NewBreaker("hello",2,3,10)
	Convey("new breaker",t, func() {
		So(b,ShouldNotBeNil)
	})
	Convey("twice new calling",t,func(){
		c := NewBreaker("hello",2,3,10)
		So(b,ShouldEqual,c)
	})
}
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
