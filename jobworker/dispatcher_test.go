package jobworker

import (
	"testing"
	"net/http"
	"time"
	"fmt"
)

func TestDispatcher_Put(t *testing.T) {
	dispatcher := NewDispather(2,10)
	//等待worker启动
	time.Sleep(time.Nanosecond * 1)
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

func TestNewHttpProxyJob(t *testing.T) {
	dispatcher := NewDispather(200,500)
	//等待worker启动
	time.Sleep(time.Nanosecond * 1)
	transport := &http.Transport{
		DisableKeepAlives : false,
		MaxIdleConnsPerHost : 10,
	}
	req,_ := http.NewRequest("GET","http://www.baidu.com",nil)
	rc := make(chan HttpProxyJobResponse)
	job := NewHttpProxyJob(transport,req,rc,"get_baidu")

	if dispatcher.Put(job){
		//accessing baidu
		select{
			//3超时控制
			case <- time.After(time.Second * 3):
				r := <- rc
				if  r.Error == nil {
					r.Response.Body.Close()
				}
				t.Log("3秒超时触发")
				job.ResetPool()
			case r := <- rc:
				if  r.Error != nil {
					t.Error(r.Error.Error())
				}else {
					resp := r.Response
					defer resp.Body.Close()
					//body, err := ioutil.ReadAll(resp.Body)
					//buf := new(bytes.Buffer)
					//buf.ReadFrom(resp.Body)
					//t.Log(buf.String())
					t.Log(fmt.Sprintf("请求百度返回 http status : %d",resp.StatusCode))

				}
				job.ResetPool()
		}
	}else{
		//queue is full
	}

}