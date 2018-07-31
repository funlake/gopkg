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
	dispatcher := NewDispather(2,10)
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
	dispatcher := NewDispather(20000,500000)
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			dispatcher.Put(&simpleJob{})
		}
	})
}
func makeRequest(dispatcher *Dispatcher) error{
	now := time.Now()
	req,_ := http.NewRequest("GET","http://www.etcchebao.com",nil)
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
				//body, err := ioutil.ReadAll(resp.Body)
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				//t.Log(fmt.Sprintf("请求百度返回 http status : %d",resp.StatusCode))
				log.Success("trasport 请求返回 http status : %d,请求时间 : %s",resp.StatusCode,time.Since(now))
			}
			job.Release()
		}
	}else{
		//queue is full
	}
	return nil
}
func makeRequestWithFastHttp(dispatcher *Dispatcher) error{
	now := time.Now()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://www.etcchebao.com")
	req.Header.SetMethod("GET")
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
				r.Response.Body()
				//body, err := ioutil.ReadAll(resp.Body)
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				//t.Log(fmt.Sprintf("请求百度返回 http status : %d",resp.StatusCode))
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
	dispatcher := NewDispather(200,500)
	makeRequestWithFastHttp(dispatcher)
}
func TestNewHttpProxyJob(t *testing.T) {
	dispatcher := NewDispather(200,500)
	makeRequest(dispatcher)
}

func BenchmarkDispatcher_BaiduWithTransport(b *testing.B) {
	b.StopTimer()
	dispatcher := NewDispather(200,500)
	b.StartTimer()
	for i:=0;i<15;i++{
		makeRequest(dispatcher)
	}
	//b.SetParallelism(10)
	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next(){
	//		makeRequest(dispatcher)
	//	}
	//})
}

//func BenchmarkDispatcher_BaiduWithFasthttp(b *testing.B) {
//	b.StopTimer()
//	dispatcher := NewDispather(200,500)
//	b.StartTimer()
//	for i:=0;i<15;i++{
//		makeRequestWithFastHttp(dispatcher)
//	}
//}

//func BenchmarkNewHttpProxyJob(b *testing.B) {
//	//b.SetParallelism(10)
//	//b.RunParallel(func(pb *testing.PB) {
//	//	for pb.Next(){
//	//		req,_ := http.NewRequest("GET","http://www.baidu.com",nil)
//	//		NewHttpProxyJob(transport,req,200,"get_baidu")
//	//	}
//	//})
//	for i:=0;i<20;i++{
//		req,_ := http.NewRequest("GET","http://www.baidu.com",nil)
//		NewHttpProxyJob(transport,req,200,"get_baidu")
//	}
//}