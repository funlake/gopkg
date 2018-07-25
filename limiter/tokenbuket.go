package limiter
import (
	"time"
	//"github.com/astaxie/beego"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"sync"
	"github.com/funlake/gopkg/utils/log"
	//"container/list"
	"github.com/funlake/gopkg/timer"
	"github.com/funlake/gopkg/utils"
)
var (
	rwLock  = new (sync.RWMutex)
	//var tb *TokenBucket
	//var syncTb sync.Once
	tokenMap = map[string]*TokenBucket{}
	tokenCache = map[string] int{}
	//tokenTick = &lib.Ticker{CallBacks:make(map[int] *list.List)}
	ticker = timer.NewTicker()
)
//var tickerSecond = map[time.Duration] *time.Ticker{}
const  (
	TB_PUTLEFT  = 1
	TB_PUTRATE = 2
)

func initTimeTokenBucket(bucketKey string,rate int,size int,second  time.Duration){
	bucketType := bucketKey[len(bucketKey) - 4 :]
	rwLock.Lock()
	if _,ok:=tokenMap[bucketKey];ok {
		//限流参数有变动，则刷新令牌桶对象
		if tokenCache[bucketKey+"_rate"] != rate ||  tokenCache[bucketKey+"_size"] != size{
			rwLock.Unlock()
			setRateSizeCache(bucketKey,rate,size)
			//参数有变,重置令牌桶
			go refreshTimeTokenBucket(bucketKey,rate,size,second)
		}else{
			rwLock.Unlock()
		}

		//桶已生成，则返回
		return
	}else{
		rwLock.Unlock()
	}
	//syncTb.Do(func() {
	// rate 大 size 小， 则为漏斗型，上大，下小 ,token消耗赶不上token放入,一般不会用到
	// rate 小 size 大 ，则为喇叭形，上小，下大 , token放入赶不上token消耗
	// rate == size,       则为直桶型 ，均匀速率，可用于长周期限流，比如一天限制请求2000次
	rwLock.Lock()
	tokenMap[bucketKey] = &TokenBucket{make(chan int ,size ), rate, size, bucketKey, false, second, time.NewTicker(time.Second * 60)}
	rwLock.Unlock()

	//初始化令牌桶,否则桶子为空，则第一次请求会被限流
	//TODO:
	// 如何记录相应桶子里,相应接入商还剩下多少token,才能在服务重启后限流进行延续，对长周期限流有意义
	//DONE
	if bucketType == "_day" {
		//因为有以下的LoopDay处理，即使重启时间在凌晨0点，如果发生跨天脏缓存问题, token也会很快被自动刷满
		tokenMap[bucketKey].putTokenIntoBucket(TB_PUTLEFT)
	}else{
		//因为短时间的token重置有跨时脏缓存问题，比如
		//1分钟限100次访问，用户在第59秒消费了99个token
		//然后服务在1分钟后重启，则新的一分钟，如果使用了载入token缓存
		//则只会载入1个token,与实际情况不相符,固短周期限流不采用重启续限机制
		tokenMap[bucketKey].putTokenIntoBucket(TB_PUTRATE)
	}
	//开routine,不停尝试往令牌桶里添加token ,满则溢出
	go tokenMap[bucketKey].loopSecond()
	//在下一个凌晨0点重启令牌桶定时植入令牌的定时器,保证接下来大体可按照固定0点到24点限流
	if bucketType == "_day" {
		go tokenMap[bucketKey].dayLimitRefreshCheck()
	}
	go tokenMap[bucketKey].catchExit()
	setRateSizeCache(bucketKey,rate,size)
	//	})
}
//获取appid对应api方法的限流令牌桶
func GetTimeTokenBucket(appid string,rate int,size int,duration  time.Duration)  *TokenBucket {
	initTimeTokenBucket(appid,rate,size,duration)
	return tokenMap[appid]
}
//刷新令牌桶
func refreshTimeTokenBucket(appid string,rate int,size int,duration  time.Duration) {
	rwLock.RLock()
	if _,ok:=tokenMap[appid];ok {
		rwLock.RUnlock()
		defer utils.RoutineRecover()
		log.Info("限流令牌tick[%d][%s],%d,%d桶重启中:",int(duration),appid,rate,size)
		tokenMap[appid].removeTokenBucket()
		initTimeTokenBucket(appid, rate, size, duration)
		log.Success("限流令牌tick[%d][%s],%d,%d桶重启完毕:",int(duration),appid,rate,size)
	}else{
		rwLock.RUnlock()
	}
}
//主结构
type TokenBucket struct {
	//令牌桶，控制TPS
	Bucket chan int
	//每固定一段时间(Dur参数)放多少个token
	Rate int
	//令牌桶大小,直接与流量相关
	BucketSize int
	//对应appid
	Appid string
	//桶是否已满
	BucketIsFull bool
	//loop second
	Second time.Duration
	////定时器,定时往桶里放令牌
	//Ticker *time.Ticker
	//每日定时器,只对每日限流起作用
	DayTicker *time.Ticker
}
//每隔一段时间执行令牌放入令牌桶动作
func (this *TokenBucket) loopSecond()  {
	//往全局定时器插入当前事件
	ticker.Set(int(this.Second),this.Appid, func() {
		this.putTokenIntoBucket(TB_PUTRATE)
	})
	log.Success("限流令牌tick[%d][%s],%d,%d桶定时器启动:",int(this.Second),this.Appid,this.Rate,this.BucketSize)
	//if tokenTick.Remove(int(this.Second),this.Appid) {
	//	tokenTick.Push(int(this.Second), this.Appid, func() {
	//		this.putTokenIntoBucket(TB_PUTRATE)
	//	})
	//	log.Success("限流令牌tick[%d][%s],%d,%d桶定时器启动:",int(this.Second),this.Appid,this.Rate,this.BucketSize)
	//}
	//for{
	//	<- this.Ticker.C
	//	this.putTokenIntoBucket(TB_PUTRATE)
	//}
}
//重启定时动作
func (this *TokenBucket) restartLoop() {
	//this.Ticker.Stop()
	log.Warning("Restarting loop second")
	if this.DayTicker != nil {
		this.DayTicker.Stop()
	}
	this.clearTokenInBucket()

	this.putTokenIntoBucket(TB_PUTRATE)
	//this.Ticker = time.NewTicker(time.Second * this.Second)
	go this.loopSecond()
}
//实现的并不优雅，此方法只针对基于每天限流的令牌桶规则
//原始问题:因为基于每天的限流规则是固定00:00 ~ 23:59，并不像
//时分秒等限流，可以做到每隔一段时间置入令牌,固目前有以下做法
//做法: 一旦网关启动后， 第一次检测到凌晨0点，就刷新令牌桶的定时任务，这样定时任务
//开始时间就会从下一个0点开始，这样就大致能做到每隔24小时刷新令牌桶
func (this *TokenBucket) dayLimitRefreshCheck() {
	//每10分钟执行一次，发现是凌晨0点则重启LoopSecond,清空令牌桶
	this.DayTicker = time.NewTicker( time.Second * 600 )
	for{
		<- this.DayTicker.C
		now := time.Now()
		if now.Hour() == 0 {
			log.Warning("Restart loop for day limit")
			//重新启动
			go this.restartLoop()
			break
		}
	}
}
//TODO:
//考虑分布式场景，考虑服务重启场景
//考虑用消息队列,可发一条指令到消息队列
//然后在回调处做植入channel的操作
func (this *TokenBucket) putTokenIntoBucket(putType int)  {
	defer utils.RoutineRecover()
	if this.BucketIsFull {
		return
	}
	rate := this.Rate
	switch putType {
	case TB_PUTLEFT:
		cache := lib.NewCache()
		r,err := cache.GetRateLimitLeft(this.Appid)
		if err == nil{
			rate,_ = strconv.Atoi(r)
			log.Info("Set rate from last stop cache:"+r)
		}
	case TB_PUTRATE:
	default:
		//rate = this.Rate
	}
	for v:=rate;v >0;v--{
		select {
		case this.Bucket <-1 :
			//log.Warning("Successfully put one token into bucket")
			continue
			//case <-time.After(time.Microsecond * 1):
		default:
			//log.Warning("Second bucket is full of tokens")
			this.BucketIsFull = true
			break
		}
	}
}
//清空令牌桶
//TODO:
//考虑分布式场景，考虑服务重启场景
func (this *TokenBucket) clearTokenInBucket()  {
	defer utils.RoutineRecover()
theEnd:
	for {
		select {
		//case <-time.After(time.Microsecond * 1):
		default:
			this.BucketIsFull = false
			log.Warning("Bucket was flushed")
			break theEnd
			//return
		case <-this.Bucket:
			//log.Warning("Get token")
			continue
		}
	}


}
//删除令牌桶
//TODO:
//考虑分布式场景，考虑服务重启场景
func (this *TokenBucket) removeTokenBucket()  {
	//if tokenTick.Remove(int(this.Second),this.Appid){
	rwLock.Lock()
	//tokenTick.Remove(int(this.Second),this.Appid)
	//tokenMap[this.Appid].Ticker.Stop()
	if tokenMap[this.Appid].DayTicker != nil {
		tokenMap[this.Appid].DayTicker.Stop()
	}
	close(tokenMap[this.Appid].Bucket)
	//tokenMap[this.Appid].clearTokenInBucket()
	delete(tokenMap, this.Appid)
	rwLock.Unlock()
	//}
}
//消费令牌
//TODO:
//考虑分布式场景，考虑服务重启场景
//考虑用消息队列,可发一条指令到消息队列
//然后在回调处做消费channel的操作
func (this *TokenBucket)ConsumeTimeToken()  bool{
	select {
	//case <- time.After(time.Millisecond * 1) :
	default:
		//log.Warning("Bucket is empty")
		return false
	case <- this.Bucket :
		//log.Warning("Get token")
		this.BucketIsFull = false
		return true
	}
	return false
}
//捕获退出，保存进程信息
func (this *TokenBucket) catchExit()  {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGINT, syscall.SIGTERM,syscall.SIGKILL)
	defer utils.RoutineRecover()
	// Block until a signal is received.
	cache := lib.NewCache()
	go func(){
		defer utils.RoutineRecover()
		select {
		case <-c :

			//fmt.Println(this.Appid+" has set left cache")
			cache.SetRateLimitLeft(this.Appid,strconv.Itoa(len(this.Bucket)))
			signal.Stop(c)
			close(c)
		}
	}()
}
func setRateSizeCache(bucketKey string,rate int,size int) {
	rwLock.Lock()
	tokenCache[bucketKey+"_rate"] = rate
	tokenCache[bucketKey+"_size"] = size
	rwLock.Unlock()
}