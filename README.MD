常用go组件，各组件用法详看各自组件目录readme或_test文件
1. timer组件
> 基础定时器是timewheel时间轮实现,所有定时任务最终只会由一个每秒定时器触发，减少多定时ticker产生的cpu资源消耗,目前定时粒度只支持至秒，分，时

2. timercache组件

> 缓存+自动定时更新，特定场景下，可利用本地缓存，减轻缓存服务器压力

3. jobworker组件

> 任务分发器，控制协程数量，量化任务执行能力，可自实现WorkerJob任务接口

4. limiter限流器组件

> 单机限流，基于令牌桶，比如300秒允许某项动作执行50次

5. circuit breaker简易熔断器
```
NewBreaker("request google",2,30,10,3)
```
> 以上规则是30秒内，超时时间超过2秒，次数超过3次,且错误比率超过10%，则启动熔断机制，等待下个30秒，自动启动半开模式，50%几率穿透，
> 30秒内假如没有超时或pass几率超过90%，则视为服务恢复

6. stats.increment 每秒统计次数器
##### demo
```
increment := &stats.Increment{}
increment.Setup("request_baidu") //添加统计
go func() {
    tm := timer.NewTimer()
    tm.Ready()
    tm.SetInterval(1, func() {
        //每秒钟输出当前qps
        log.Info("qps : %d",increment.GetStat("request_baidu").GetReport().GetQps())
    })
}()

for  {
    increment.IncrRequest("request_baidu")
    time.Sleep(time.Millisecond * 300) //当于1秒请求3,4次
}
```


