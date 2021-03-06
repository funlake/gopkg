package limiter

import (
	"github.com/funlake/gopkg/timer"
	"github.com/funlake/gopkg/utils"
	"github.com/funlake/gopkg/utils/log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var tbTicker = timer.NewTicker()

const (
	TB_PUTLEFT = 1 // 每日限流标识
	TB_PUTRATE = 2 //  秒，分，时限流标识
)

func NewTokenBucketSchedular() *tokenBucketSchedular {
	return &tokenBucketSchedular{}
}

type tokenBucketSchedular struct {
	tbKv sync.Map
	tcKv sync.Map
}

func (tbs *tokenBucketSchedular) GetTimeTokenBucket(bucketKey string, rate int, qps int, second time.Duration, dayRateFun *DayRateFun) *tokenBucket {
	return tbs.makeTimeBucket(bucketKey, rate, qps, second, dayRateFun)
}
func (tbs *tokenBucketSchedular) makeTimeBucket(bucketKey string, rate int, size int, second time.Duration, dayRateFun *DayRateFun) *tokenBucket {
	//tbs.Lock()
	//defer tbs.Unlock()
	if tbc, ok := tbs.tbKv.Load(bucketKey); ok {
		crate, ok2 := tbs.tcKv.Load(bucketKey + "_rate")
		csize, ok3 := tbs.tcKv.Load(bucketKey + "_size")
		if ok2 && ok3 && (crate.(int) != rate || csize != size) {
			tbs.setRateSize(bucketKey, rate, size)
			tbs.restartTimeBucket(bucketKey, rate, size, second, dayRateFun)
		}
		//tbs.Unlock()
		return tbc.(*tokenBucket)
	}
	//tbs.tbHash[bucketKey] = &tokenBucket{bucket:make(chan int ,size ), rate:rate,size:size, bucketKey:bucketKey, bucketIsFull:false, second:second, ticker:tbTicker,dayTicker:true}
	tb := &tokenBucket{bucket: make(chan int, size), rate: rate, size: size, bucketKey: bucketKey, bucketIsFull: false, second: second, ticker: tbTicker, dayTicker: true}
	tbs.tbKv.Store(bucketKey, tb)
	//tbs.Unlock()
	bucketType := bucketKey[len(bucketKey)-4:]
	rateType := TB_PUTRATE
	if bucketType == "_day" {
		rateType = TB_PUTLEFT
	}
	if dayRateFun != nil {
		//tbs.tbHash[bucketKey].setDayRateFun(dayRateFun)
		tb.setDayRateFun(dayRateFun)
	}
	//初始化令牌数
	tb.putTokenIntoBucket(rateType)
	//开启定时令牌刷新
	tb.startTimer()
	//缓存配置
	tbs.setRateSize(bucketKey, rate, size)
	//获取信号,缓存每日限流数据
	utils.WrapGo(func() {
		tb.catchExit()
	}, "token bucket cache exit")
	return tb
}
func (tbs *tokenBucketSchedular) restartTimeBucket(bucketKey string, rate int, size int, duration time.Duration, dayRateFun *DayRateFun) {
	if oldTb, ok := tbs.tbKv.Load(bucketKey); ok {
		log.Info("限流令牌tick[%d][%s],%d,%d桶重启中:", int(duration), bucketKey, rate, size)
		oldTb.(*tokenBucket).stopTokenBucketTimer()
		//delete(tbs.tbHash, bucketKey)
		tbs.tbKv.Delete(bucketKey)
		//tbs.Unlock()
		tbs.makeTimeBucket(bucketKey, rate, size, duration, dayRateFun)
		log.Success("限流令牌tick[%d][%s],%d,%d桶重启完毕:", int(duration), bucketKey, rate, size)
	}
}
func (tbs *tokenBucketSchedular) setRateSize(bucketKey string, rate int, size int) {
	//tbs.Lock()
	//tbs.tbCache[bucketKey+"_rate"] = rate
	//tbs.tbCache[bucketKey+"_size"] = size
	tbs.tcKv.Store(bucketKey+"_rate", rate)
	tbs.tcKv.Store(bucketKey+"_size", size)
	//tbs.Unlock()
}

type DayRateFun struct {
	Set func(key string, val string)
	Get func(key string) (string, error)
}

//主结构
type tokenBucket struct {
	//令牌桶，控制TPS
	bucket chan int
	//每固定一段时间(Dur参数)放多少个token
	rate int
	//令牌桶大小,直接与流量相关
	size int
	//对应appid
	bucketKey string
	//桶是否已满
	bucketIsFull bool
	//loop second
	second time.Duration
	//定时器,定时往桶里放令牌
	ticker *timer.Ticker
	//每日定时器,只对每日限流起作用
	dayTicker bool
	//日限流缓存器
	dayFun *DayRateFun
}

func (tbk *tokenBucket) setDayRateFun(fs *DayRateFun) {
	log.Warning("Set day rate limit fun")
	tbk.dayFun = fs
}

//每隔一段时间执行令牌放入令牌桶动作
func (tbk *tokenBucket) startTimer() {
	//往全局定时器插入当前事件
	tbk.ticker.Set(int(tbk.second), tbk.bucketKey, func() {
		tbk.putTokenIntoBucket(TB_PUTRATE)
	})
	//log.Success("限流令牌桶%s,%dreq/%ds,rate:%d桶定时器启动:",tbk.bucketKey,tbk.size,int(tbk.second),tbk.rate)
}

//重启定时动作
func (tbk *tokenBucket) restartTimer() {
	//tbk.Ticker.Stop()
	log.Warning("Restarting loop second")
	if tbk.dayTicker {
		//tbk.dayTicker.Stop()
		tbk.ticker.Stop(600, tbk.bucketKey+"_600")
	}
	tbk.clearTokenInBucket()
	tbk.putTokenIntoBucket(TB_PUTRATE)
	//tbk.Ticker = time.NewTicker(time.second * tbk.second)
	tbk.startTimer()
}

//实现的并不优雅，此方法只针对基于每天限流的令牌桶规则
//原始问题:因为基于每天的限流规则是固定00:00 ~ 23:59，并不像
//时分秒等限流，可以做到每隔一段时间置入令牌,固目前有以下做法
//做法: 一旦网关启动后， 第一次检测到凌晨0点，就刷新令牌桶的定时任务，这样定时任务
//开始时间就会从下一个0点开始，这样就大致能做到每隔24小时刷新令牌桶
func (tbk *tokenBucket) dayLimitRefreshCheck() {
	//每10分钟执行一次，发现是凌晨0点则重启LoopSecond,清空令牌桶
	tbk.ticker.Set(600, tbk.bucketKey+"_600", func() {
		log.Info("day loop checking")
		now := time.Now()
		if now.Hour() == 0 {
			log.Warning("Restart loop for day limit")
			//重新启动
			tbk.restartTimer()
			tbk.ticker.Stop(600, tbk.bucketKey+"_600")
		}
	})
}

//TODO:
//考虑分布式场景，考虑服务重启场景
//考虑用消息队列,可发一条指令到消息队列
//然后在回调处做植入channel的操作
func (tbk *tokenBucket) putTokenIntoBucket(putType int) {
	//defer utils.RoutineRecover()
	if tbk.bucketIsFull {
		return
	}
	rate := tbk.rate
	switch putType {
	case TB_PUTLEFT:
		//cache := lib.NewCache()
		//r,err := cache.GetRateLimitLeft(tbk.bucketKey)
		if tbk.dayFun != nil {
			r, err := tbk.dayFun.Get(tbk.bucketKey)
			if err == nil {
				rate, _ = strconv.Atoi(r)
				log.Info("Set day rate from last stop:" + r)
			}
		}
	case TB_PUTRATE:
	default:
		//rate = tbk.rate
	}
	for v := rate; v > 0; v-- {
		select {
		case tbk.bucket <- 1:
			//log.Warning("Successfully put one token into bucket")
			continue
			//case <-time.After(time.Microsecond * 1):
		default:
			//log.Warning("second bucket is full of tokens")
			tbk.bucketIsFull = true
			break
		}
	}
}

//清空令牌桶
//TODO:
//考虑分布式场景，考虑服务重启场景
func (tbk *tokenBucket) clearTokenInBucket() {
theEnd:
	for {
		select {
		//case <-time.After(time.Microsecond * 1):
		default:
			tbk.bucketIsFull = false
			log.Warning("bucket was flushed")
			break theEnd
			//return
		case <-tbk.bucket:
			//log.Warning("Get token")
			continue
		}
	}
}
func (tbk *tokenBucket) stopTokenBucketTimer() {
	if tbk.dayTicker {
		tbk.ticker.Stop(600, tbk.bucketKey+"_600")
	}
	//停掉定时器,否则定时器还在，channel却已关闭
	tbk.ticker.Stop(int(tbk.second), tbk.bucketKey)
	close(tbk.bucket)
}

//消费令牌
//TODO:
//考虑分布式场景，考虑服务重启场景
//考虑用消息队列,可发一条指令到消息队列
//然后在回调处做消费channel的操作
func (tbk *tokenBucket) GetToken() bool {
	select {
	case <-tbk.bucket:
		tbk.bucketIsFull = false
		return true
	default:
	}
	return false
}

//捕获退出，保存进程信息
func (tbk *tokenBucket) catchExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	select {
	case <-c:
		log.Warning(tbk.bucketKey + " caught stop signal")
		if tbk.dayFun != nil {
			tbk.dayFun.Set(tbk.bucketKey, strconv.Itoa(len(tbk.bucket)))
		}
		signal.Stop(c)
		close(c)
	}
}