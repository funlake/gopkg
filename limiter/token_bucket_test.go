package limiter

import (
	"testing"
	"golang.org/x/time/rate"
	"time"
	"context"
)

func Benchmark_GetTimeTokenBucket100(b *testing.B) {

	//runtime.GOMAXPROCS(1)
	//for i:=0;i<b.N;i++{
	//	tb := tokenBucketSchedular.GetTimeTokenBucket("test_second", 100, 100, 1,nil)
	//	tb.GetToken()
	//}
	tokenBucketSchedular := NewTokenBucketSchedular()
	b.SetParallelism(20)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			tokenBucketSchedular.GetTimeTokenBucket("test_second", 300, 300, 1,nil).GetToken()
		}
	})
}
//check raw rate controller implemented by go itself
//a little bit slower then mine,:p
//but raw one got Wait method to block request let it might could be executed later
//this is one thing that i didn't have it.
func BenchmarkRawRateLimit(b *testing.B) {
	b.SetParallelism(20)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			limiter := rate.NewLimiter(rate.Every(time.Second * 1),300)
			limiter.Wait(context.Background())
		}
	})
}

//func Benchmark_GetTimeTokenBucket1000(b *testing.B) {
//	tokenBucketSchedular := NewTokenBucketSchedular()
//	b.SetParallelism(20)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next(){
//			tokenBucketSchedular.GetTimeTokenBucket("test_second", 1000, 1000, 1,nil).GetToken()
//		}
//	})
//}
//func Benchmark_GetTimeTokenBucket10000(b *testing.B) {
//	tokenBucketSchedular := NewTokenBucketSchedular()
//	b.SetParallelism(20)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next(){
//			tokenBucketSchedular.GetTimeTokenBucket("test_second", 10000, 10000, 1,nil).GetToken()
//		}
//	})
//}
//func BenchmarkTokenBuckt_ConsumeTimeToken(b *testing.B) {
//	tokenBucketSchedular := NewTokenBucketSchedular()
//	rl := 1000
//	//10秒1000次请求
//	secondBucket := tokenBucketSchedular.GetTimeTokenBucket("test_second", rl, rl, 10,nil)
//	//模拟1秒并发rl次
//	for i:=0;i<rl;i++{
//		secondBucket.GetToken()
//	}
//	b.SetParallelism(20)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next(){
//			if secondBucket.GetToken(){
//				b.Error("Your are under attack!")
//			}
//		}
//	})
//}
//func TestRateLimit(t *testing.T) {
//	tokenBucketSchedular := NewTokenBucketSchedular()
//	rl := 1000
//	secondBucket := tokenBucketSchedular.GetTimeTokenBucket("test_second", rl, rl, 1,nil)
//	//模拟1秒并发rl次
//	for i:=0;i<rl;i++{
//		secondBucket.GetToken()
//	}
//	//同一秒第rl+1次请求进来
//	if secondBucket.GetToken(){
//		t.Error("Token Consume fail")
//	}else{
//		t.Log(fmt.Sprintf("TokenBucket with %d second,%d size stop your request successfully",secondBucket.second,rl))
//	}
//}

func TestTokenBucket_GetToken(t *testing.T) {
	tokenBucketSchedular := NewTokenBucketSchedular()
	for i:=0;i<1000;i++{
		tokenBucketSchedular.GetTimeTokenBucket("test_second", 100, 1001, 1,nil).GetToken()
	}
}

