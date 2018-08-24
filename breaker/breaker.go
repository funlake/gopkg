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
	b := &breaker{id:id,rate:rate,status:0,errChans:make(chan breakerItem,100)}
	b.init(id,timeout,window)
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
	metrics *MestricsEntity
}

func (b *breaker) init(id string,timeout int,window int){
	b.metrics = NewMetrics().NewEntity(id,timeout, window)
	go b.tick()
}
func (b *breaker) Run(fun func (),okfun func(),failfun func(run bool)){
	run := true
	if b.isOpen(){
		if len(b.errChans) > 0{
			log.Warning("%d,%d,%d",b.metrics.pass,len(b.errChans), len(b.errChans) / (b.metrics.pass + len(b.errChans)))
			if ( len(b.errChans) / (b.metrics.pass + len(b.errChans)) ) * 100 >= b.rate {
				b.status = 2
				run = false
			}
		}else{
			if b.status > 0 {
				b.status = 0
			}
		}
	}
	if b.isClose(){
		go failfun(false)
		return
	}
	if b.isHalfopen(){
		if b.metrics.pass > 0 {
			if ( b.metrics.pass / (b.metrics.pass + len(b.errChans)) ) * 100 >= (100 - b.rate) {
				b.status = 0
				run = true
			}
		}
		if !b.isOpen() {
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(100)
			if i > 50 {
				run = true
			} else {
				run = false
				go failfun(false)
				return
			}
		}
	}
	if run {
		cxt, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(b.metrics.GetTimeout()))

		ch := make(chan bool)
		go func() {
			fun()
			ch <- true
		}()
		select {
		case <-cxt.Done():
			log.Success("触发超时:%d秒",b.metrics.GetTimeout())
			b.metrics.broken = b.metrics.broken + 1
			if b.status == 1{
				b.status = 2
			}
			utils.WrapGo(func() {
				go failfun(true)
				b.errChans <- breakerItem{
					notations:b.id,
				}
			},"breaking")
			return
		case <-ch:
			b.metrics.pass = b.metrics.pass + 1
			go okfun()
			return
		}
	}
}

func (b *breaker) tick(){
	breakerTimer.SetInterval(b.metrics.GetWindow(), func() {
		b.metrics.pass = 0
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
	breakerTimer.SetInterval(b.metrics.GetWindow() * 2, func() {
		if b.status == 2 {
			b.status = 1
		}
	})
}

func (b *breaker) isHalfopen() bool{
	return b.status == 1
}

func (b *breaker) isClose() bool{
	return b.status == 2
}

func (b *breaker) isOpen() bool{
	return b.status == 0
}