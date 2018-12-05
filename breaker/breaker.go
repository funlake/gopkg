package breaker

import (
	"context"
	"time"
	"github.com/funlake/gopkg/timer"
	"sync"
	"math/rand"
	//"github.com/funlake/gopkg/utils"
	"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/utils"
	"fmt"
)
var breakerTimer = timer.NewTimer()
var breakerMap sync.Map
func init(){
	breakerTimer.Ready()
}
func NewBreaker(id string) *breaker {
	if c,ok := breakerMap.Load(id);ok{
		return c.(*breaker)
	}
	//b := &breaker{id:id,rate:rate,status:0,/*errChans:make(chan breakerItem,100),*/timeout:timeout,window:window,pass:0,broken:0,min:min}
	b := &breaker{id:id,pass:0,broken:0}
	b.init()
	breakerMap.Store(id,b)
	return b
}
type breaker struct{
	id string
	rate int
	status int
	//errChans chan breakerItem
	//metrics *MestricsEntity
	timeout int
	pass int
	broken int
	window int
	//至少出现多少次超时才熔断
	min int
	//half open ticker
	hot *timer.SlotItem
}

func (b *breaker) init(){
	//b.metrics = NewMetrics().NewEntity(id,timeout, window)
}
func (b *breaker) Run(fun func (),okfun func(),failfun func(run bool)){
	run := true
	//if b.pass+b.broken > 0 {
	//	log.Warning("超时次数：%d ，正常次数：%d ,比率 ：%d",b.broken,b.pass,(b.broken * 100 / (b.broken + b.pass)) )
	//}
	if b.isClose(){
		if(b.broken >= b.min) {
			if (b.broken * 100 / (b.pass + b.broken)) >= b.rate{
				log.Error("%s熔断开启,超时请求比率: %d%%",b.id,(b.broken * 100 / (b.pass + b.broken) ))
				b.open()
				b.tick()
			}
		}else{
			if !b.isClose() {
				b.close()
				b.untick()
			}
		}
	}
	if b.isOpen(){
		utils.WrapGo(func() {failfun(false)},"failfun-open")
		return
	}
	if b.isHalfopen(){
		if b.pass > 0 {
			//if ( b.pass / (b.pass + len(b.errChans)) ) * 100 >= (100 - b.rate) {
			if (b.pass * 100 / (b.pass + b.broken)) >= (100 - b.rate){
				log.Info(fmt.Sprintf("%s熔断关闭",b.id))
				b.close()
				run = true
			}
		}
		if !b.isClose() {
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(100)
			if i > 50 {
				run = true
			} else {
				utils.WrapGo(func() {failfun(false)},"failfun-halfopen")
				return
			}
		}
	}
	if run {
		cxt, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(b.timeout))
		ch := make(chan struct{})
		//mention: fun() can not be a blocking behaviour
		utils.WrapGo(func() {
			fun()
			ch <- struct{}{}
		},"breaker_run")
		select {
			case <-cxt.Done():
				if b.isHalfopen(){
					b.open()
				}
				b.broken = b.broken + 1
				utils.WrapGo(func() {
					failfun(true)
					<-ch
				},"failfun")
				return
			case <-ch:
				b.pass = b.pass + 1
				utils.WrapGo(func() {
					okfun()
				},"okfun")
				return
		}
	}
}

func (b *breaker) tick(){
	b.hot = breakerTimer.SetInterval(b.window, func() {
		b.pass = 0
		b.broken = 0
		if b.isOpen() {
			log.Info(fmt.Sprintf("%s熔断半开",b.id))
			b.halfopen()
		}
	})
}

func (b *breaker) untick(){
	breakerTimer.StopInterval(b.hot)
}
func (b *breaker) SetWindow(window int)  {
	b.window = window
}

func (b *breaker) SetRate(rate int){
	b.rate = rate
}
func (b *breaker) SetMin(min int){
	b.min = min
}
func (b *breaker) SetTimemout(timeout int){
	b.timeout = timeout
}
func (b *breaker) isHalfopen() bool{
	return b.status == 1
}

func (b *breaker) halfopen(){
	b.status = 1
}
func (b *breaker) isClose() bool{
	return b.status == 0
}
func (b *breaker) close()  {
	b.status = 0
}
func (b *breaker) isOpen() bool{
	return b.status == 2
}
//打开熔断器
func (b *breaker) open(){
	b.status = 2
}
func (b *breaker) IsBroken() bool{
	return b.status == 2
}
