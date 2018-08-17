package breaker

import (
	"context"
	"time"
	"github.com/funlake/gopkg/timer"
	"sync"
	"math/rand"
	"github.com/funlake/gopkg/utils"
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
func (b *breaker) Run(fun func (),okfun func(),failfun func()){
	run := true
	if b.isClose(){
		go failfun()
		return
	}
	if b.isHalfopen(){
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(100)
		if i > 50{
			run = true
		}else{
			run = false
			go failfun()
			return
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
			b.metrics.broken = b.metrics.broken + 1
			if b.status == 1{
				b.status = 2
			}
			utils.WrapGo(func() {
				failfun()
				b.errChans <- breakerItem{
					notations:b.id,
				}
			},"breaking")
		case <-ch:
			b.metrics.pass = b.metrics.pass + 1
			go okfun()
		}
	}
}

func (b *breaker) tick(){
	breakerTimer.SetInterval(b.metrics.GetWindow(), func() {
		if len(b.errChans) > 0{
			if ( len(b.errChans) / (b.metrics.pass + len(b.errChans)) ) * 100 >= b.rate {
				b.status = 2
			}
		}else{
			if b.status > 0 {
				b.status = 0
			}
		}
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