package breaker

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

//func TestBreaker_Run(t *testing.T) {
//	bre := NewBreaker("hello",2,3,10)
//	for {
//		bre.Run(func() {
//			res,_ := http.Get("http://admin.etcchebao.com/index2.php")
//			log.Success("%d",res.StatusCode)
//			//now := time.Now()
//			//rand.Seed(time.Now().UnixNano())
//			//t := rand.Intn(5)
//			//time.Sleep(time.Second * time.Duration(t))
//			//log.Info("%s",time.Since(now))
//		}, func() {
//			log.Success("goog job")
//		}, func() {
//			log.Error("break!")
//		})
//		time.Sleep(time.Millisecond * 100)
//	}
//}
func BenchmarkBreaker_Run(b *testing.B) {
	b.SetParallelism(20)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			NewBreaker("hello",2,3,10)
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
