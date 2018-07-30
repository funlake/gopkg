##### Intro
单机版令牌桶限速器
###### Install
go get github.com/funlake/gopkg/limiter
###### Import
import github.com/funlake/gopkg/limiter
```
tokenBucketSchedular := limiter.NewTokenBucketSchedular()
// 100/s qps limiter
//第一个100 是每1秒令牌置入的速度
//第二个100 是qps
rl := tokenBucketSchedular.GetTimeTokenBucket("test_second", 100, 100, 1,nil)
if rl.GetToken(){
    //pass
}else{
    //deny
}
// 200/100s , 每100秒允许200次请求
rl2 := tokenBucketSchedular.GetTimeTokenBucket("test_second", 100, 200, 100,nil)
if rl2.GetToken(){
    //pass
}else{
    //deny
}
```