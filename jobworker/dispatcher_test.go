package jobworker

import (
	"fmt"
	"github.com/funlake/gopkg/utils/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
	"time"
)

var transport = &http.Transport{
	DisableKeepAlives:   false,
	MaxIdleConnsPerHost: 10,
}
var fasthttpClient = &fasthttp.Client{}

func TestDispatcher_Put(t *testing.T) {
	dispatcher := NewBlockingDispather(2, 10)
	job := NewSimpleJob(func() {
		//do nothing
	}, func() string {
		return "simplejob"
	}, func(dispatcher *BlockingDispatcher) {
		log.Error("Worker exhaust!")
	})
	for i := 0; i < 11; i++ {
		dispatcher.Put(job)
	}
	if dispatcher.Put(job) {
		t.Log("Ok")
	} else {
		t.Log("job queue is full")
	}
}

func BenchmarkDispatcher_Put(b *testing.B) {
	dispatcher := NewBlockingDispather(20000, 500000)
	job := NewSimpleJob(func() {
		//do nothing
	}, func() string {
		return "simplejob"
	}, func(dispatcher *BlockingDispatcher) {
		log.Error("Worker exhaust!")
	})
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dispatcher.Put(job)
		}
	})
}

func makeRequest(dispatcher *NonBlockingDispatcher) error {
	now := time.Now()
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	job := NewHttpProxyJob(transport, req, "get_baidu")
	if dispatcher.Put(job) {
		//accessing baidu
		select {
		//3s超时控制
		case <-time.After(time.Second * 3):
			r := <-job.GetResChan()
			if r.Error == nil {
				r.Response.Body.Close()
			}
		case r := <-job.GetResChan():
			if r.Error != nil {
				return r.Error
			} else {
				resp := r.Response
				defer resp.Body.Close()
				//返回结果
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				log.Success("trasport 请求返回 http status : %d,请求时间 : %s", resp.StatusCode, time.Since(now))
			}
		}
	} else {
		log.Warning("disptacher is closed or full of tasks")
		//queue is full
	}
	return nil
}

func makeRequest2(dispatcher *BlockingDispatcher) error {
	now := time.Now()
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	job := NewHttpProxyJob(transport, req, "get_baidu")
	if dispatcher.Put(job) {
		//accessing baidu
		select {
		//3s超时控制
		case <-time.After(time.Second * 3):
			r := <-job.GetResChan()
			if r.Error == nil {
				r.Response.Body.Close()
			}
		case r := <-job.GetResChan():
			if r.Error != nil {
				return r.Error
			} else {
				resp := r.Response
				defer resp.Body.Close()
				//返回结果
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				log.Success("trasport 请求返回 http status : %d,请求时间 : %s", resp.StatusCode, time.Since(now))
			}
		}
	} else {
		//queue is full
	}
	return nil
}

func makeRequestWithFastHttp(dispatcher *NonBlockingDispatcher) error {
	url := "http://www.baidu.com"
	now := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	//req.AppendBodyString(qs)
	job := NewFastHttpProxyJob(fasthttpClient, req, "get_baidu_withfasthttp")
	if dispatcher.Put(job) {
		//accessing baidu
		select {
		//3s超时控制
		case <-time.After(time.Second * 3):
			r := <-job.GetResChan()
			if r.Error == nil {
				r.Response.Body()
			}

		case r := <-job.GetResChan():
			if r.Error != nil {
				return r.Error
			} else {
				//r.Response.Body()
				log.Success("fasthttp 请求返回 http status : %d,请求时间 : %s", r.Response.StatusCode(), time.Since(now))
			}
		}
	} else {
		//queue is full
	}
	return nil
}
func makeSimpleBlockingJob(dispatcher *BlockingDispatcher) error{
	job := NewSimpleJob(func() {
		fmt.Printf("great world!")
		time.Sleep(time.Second * 2)
	}, func() string {
		return "simplejob"
	}, func(dispatcher *BlockingDispatcher) {
		log.Error("Worker exhaust!")
	})
	dispatcher.Put(job)
	return nil
}
func makeRequestWithBlockingFastHttp(dispatcher *BlockingDispatcher) error {
	url := "http://www.baidu.com"
	now := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	//req.AppendBodyString(qs)
	job := NewFastHttpProxyJob(fasthttpClient, req, "get_baidu_withfasthttp")
	if dispatcher.Put(job) {
		//accessing baidu
		select {
		//3s超时控制
		case <-time.After(time.Second * 3):
			r := <-job.GetResChan()
			if r.Error == nil {
				r.Response.Body()
			}
		case r := <-job.GetResChan():
			if r.Error != nil {
				return r.Error
			} else {
				//r.Response.Body()
				log.Success("fasthttp 请求返回 http status : %d,请求时间 : %s", r.Response.StatusCode(), time.Since(now))
			}
		}
	} else {
		//queue is full
	}
	return nil
}
//func TestNewFastHttpProxyJob(t *testing.T) {
//	dispatcher := NewNonBlockingDispather(2, 500)
//	for i := 0; i < 50; i++ {
//		makeRequestWithFastHttp(dispatcher)
//	}
//}
//func TestNewHttpProxyJob(t *testing.T) {
//	dispatcher := NewNonBlockingDispather(2, 500)
//	makeRequest(dispatcher)
//}
func TestBlockingNewHttpProxyJob(t *testing.T) {
	dispatcher := NewBlockingDispather(1, 10)
	for i := 0; i < 500; i++ {
		makeSimpleBlockingJob(dispatcher)
		//makeRequestWithBlockingFastHttp(dispatcher)
	}
	time.Sleep(time.Second * 100)
}
//func BenchmarkDispatcher_WithTransport(b *testing.B) {
//	dispatcher := NewBlockingDispather(1, 1)
//	//b.SetParallelism(10)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			makeRequest2(dispatcher)
//		}
//	})
//}

//func TestStopDispatcher(t *testing.T) {
//	dispatcher := NewNonBlockingDispather(10, 100)
//	go func() {
//		time.Sleep(time.Second * 1)
//		log.Info("关闭dispatcher")
//		dispatcher.Stop()
//	}()
//	for i := 0; i < 100; i++ {
//		makeRequest(dispatcher)
//		time.Sleep(time.Millisecond * 200)
//	}
//}
