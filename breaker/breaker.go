package breaker

import (
	"context"
	"time"
	"github.com/funlake/gopkg/timer"
	"sync"
	"math/rand"
	"github.com/funlake/gopkg/utils"
	"github.com/funlake/gopkg/utils/log"
)
var breakerTimer = timer.NewTimer()
var breakerMap sync.Map
func init(){
	breakerTimer.Ready()
}
func NewBreaker(id string,timeout int,window int,rate int) *breaker {
	if c,ok := breakerMap.Load(id);ok{
		return c.(*breaker)
	}
	b := &breaker{id:id,rate:rate,status:0,errChans:make(chan breakerItem,100),timeout:timeout,window:window,pass:0}
	b.init()
	breakerMap.Store(id,b)
	return b
}
type breakerItem struct{
	notations string
}
type breaker struct{
	id string
	rate int
	status int
	errChans chan breakerItem
	//metrics *MestricsEntity
	timeout int
	pass int
	window int
}

func (b *breaker) init(){
	//b.metrics = NewMetrics().NewEntity(id,timeout, window)
	go b.tick()
}
func (b *breaker) Run(fun func (),okfun func(),failfun func(run bool)){
	run := true
	if b.isOpen(){
		if len(b.errChans) > 0{
			//log.Warning("%d,%d,%d",b.metrics.pass,len(b.errChans), len(b.errChans) / (b.metrics.pass + len(b.errChans)))
			if ( len(b.errChans) / (b.pass + len(b.errChans)) ) * 100 >= b.rate {
				log.Error("%s 触发熔断,超时请求比率: %d%",b.id,(len(b.errChans) / (b.pass + len(b.errChans)) )* 100)
				b.close()
				return
			}
		}else{
			if !b.isOpen() {
				b.open()
			}
		}
	}
	if b.isClose(){
		go failfun(false)
		return
	}
	if b.isHalfopen(){
		if b.pass > 0 {
			if ( b.pass / (b.pass + len(b.errChans)) ) * 100 >= (100 - b.rate) {
				b.open()
				run = true
			}
		}
		if !b.isOpen() {
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(100)
			if i > 50 {
				run = true
			} else {
				go failfun(false)
				return
			}
		}
	}
	if run {
		cxt, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(b.timeout))
		ch := make(chan bool)
		go func() {
			fun()
			ch <- true
		}()
		select {
			case <-cxt.Done():
				if b.isHalfopen(){
					b.close()
				}
				utils.WrapGo(func() {
					select{
						case b.errChans <- breakerItem{notations:b.id}:
						default:
							//full of err chan means all the things go wrong badly
							b.close()
					}
					//go failfun(true)
					//b.errChans <- breakerItem{notations:b.id}
				},"breaking")
				return
			case <-ch:
				b.pass = b.pass + 1
				go okfun()
				return
		}
	}
}

func (b *breaker) tick(){
	breakerTimer.SetInterval(b.window, func() {
		b.pass = 0
		if b.isClose() {
			b.halfopen()
		}
		go func() {
			for {
				select {
				case <- b.errChans:

				default:
					return
				}
			}
		}()

	})
}

func (b *breaker) isHalfopen() bool{
	return b.status == 1
}

func (b *breaker) halfopen(){
	b.status = 1
}
func (b *breaker) isClose() bool{
	return b.status == 2
}
func (b *breaker) close()  {
	b.status = 2
}
func (b *breaker) isOpen() bool{
	return b.status == 0
}
func (b *breaker) open(){
	b.status = 0
}
