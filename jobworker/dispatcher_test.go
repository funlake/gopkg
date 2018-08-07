package jobworker

import (
	"testing"
	"net/http"
	"time"
	"github.com/valyala/fasthttp"
	"github.com/funlake/gopkg/utils/log"
)
var transport = &http.Transport{
	DisableKeepAlives : false,
	MaxIdleConnsPerHost : 10,
}
var fasthttpClient = &fasthttp.Client{}
func TestDispatcher_Put(t *testing.T) {
	dispatcher := NewBlockingDispather(2,10)
	for i:=0;i<10;i++{
		dispatcher.Put(&simpleJob{})
	}
	if dispatcher.Put(&simpleJob{}){
		t.Log("job queue is full")
	}else{
		t.Log("Ok")
	}
}

func BenchmarkDispatcher_Put(b *testing.B) {
	dispatcher := NewBlockingDispather(20000,500000)
	job := NewSimpleJob()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			dispatcher.Put(job)
		}
	})
}
func makeRequest(dispatcher *NonBlockingDispatcher) error{
	now := time.Now()
	req,_ := http.NewRequest("GET","http://www.baidu.com",nil)
    job := NewHttpProxyJob(transport,req,200,"get_baidu")
	if dispatcher.Put(job){
		//accessing baidu
		select{
		//3s超时控制
		case <- time.After(time.Second * 3):
			r := <- job.GetResChan()
			if  r.Error == nil {
				r.Response.Body.Close()
			}
			job.Release()
		case r := <- job.GetResChan():
			if  r.Error != nil {
				return r.Error
			}else {
				resp := r.Response
				defer resp.Body.Close()
				//返回结果
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				log.Success("trasport 请求返回 http status : %d,请求时间 : %s",resp.StatusCode,time.Since(now))
			}
			job.Release()
		}
	}else{
		//queue is full
	}
	return nil
}
func makeRequestWithFastHttp(dispatcher *NonBlockingDispatcher) error{
	url := "http://www.baidu.com"
	now := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	//req.AppendBodyString(qs)
	job := NewFastHttpProxyJob(fasthttpClient,req,200,"get_baidu_withfasthttp")
	if dispatcher.Put(job){
		//accessing baidu
		select{
		//3s超时控制
		case <- time.After(time.Second * 3):
			r := <- job.GetResChan()
			if  r.Error == nil {
				r.Response.Body()
			}
			job.Release()
		case r := <- job.GetResChan():
			if  r.Error != nil {
				return r.Error
			}else {
				//r.Response.Body()
				log.Success("fasthttp 请求返回 http status : %d,请求时间 : %s",r.Response.StatusCode(),time.Since(now))
			}
			job.Release()
		}
	}else{
		//queue is full
	}
	return nil
}


func makeRequestWithBlockingFastHttp(dispatcher *BlockingDispatcher) error{
	url := "http://www.baidu.com"
	now := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	//req.AppendBodyString(qs)
	job := NewFastHttpProxyJob(fasthttpClient,req,200,"get_baidu_withfasthttp")
	if dispatcher.Put(job){
		//accessing baidu
		select{
		//3s超时控制
		case <- time.After(time.Second * 3):
			r := <- job.GetResChan()
			if  r.Error == nil {
				r.Response.Body()
			}
			job.Release()
		case r := <- job.GetResChan():
			if  r.Error != nil {
				return r.Error
			}else {
				//r.Response.Body()
				log.Success("fasthttp 请求返回 http status : %d,请求时间 : %s",r.Response.StatusCode(),time.Since(now))
			}
			job.Release()
		}
	}else{
		//queue is full
	}
	return nil
}
func TestNewFastHttpProxyJob(t *testing.T) {
	dispatcher := NewNonBlockingDispather(2,500)
	for i:=0;i<50;i++ {
		makeRequestWithFastHttp(dispatcher)
	}
}
func TestNewHttpProxyJob(t *testing.T) {
	dispatcher := NewNonBlockingDispather(2,500)
	makeRequest(dispatcher)
}
func TestBlockingNewHttpProxyJob(t *testing.T) {
	dispatcher := NewBlockingDispather(2,5)
	for i:=0;i<50;i++ {
		makeRequestWithBlockingFastHttp(dispatcher)
	}
}
//func BenchmarkDispatcher_WithTransport(b *testing.B) {
//	//b.SetParallelism(10)
//	//b.RunParallel(func(pb *testing.PB) {
//	//	for pb.Next(){
//	//		makeRequest(dispatcher)
//	//	}
//	//})
//}
