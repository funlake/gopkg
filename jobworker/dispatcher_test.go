package jobworker

import (
	"testing"
	"net/http"
	"time"
	"github.com/funlake/gopkg/utils/log"
)

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
func makeBaiduRequest(dispatcher *Dispatcher) error{
	transport := &http.Transport{
		DisableKeepAlives : false,
		MaxIdleConnsPerHost : 10,
	}
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
				//body, err := ioutil.ReadAll(resp.Body)
				//buf := new(bytes.Buffer)
				//buf.ReadFrom(resp.Body)
				//t.Log(buf.String())
				//t.Log(fmt.Sprintf("请求百度返回 http status : %d",resp.StatusCode))
				log.Success("请求百度返回 http status : %d",resp.StatusCode)
			}
			job.Release()
		}
	}else{
		//queue is full
	}
	return nil
}
func TestNewHttpProxyJob(t *testing.T) {
	dispatcher := NewDispather(200,500)
	makeBaiduRequest(dispatcher)
}

func BenchmarkDispatcher_Baidu(b *testing.B) {
	dispatcher := NewDispather(200,500)
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			makeBaiduRequest(dispatcher)
		}
	})
}