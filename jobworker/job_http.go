package jobworker

import (
	"errors"
	"github.com/funlake/gopkg/utils/log"
	"net/http"
	"sync"
	"time"
)

type HttpProxyJobResponse struct {
	Response *http.Response
	Error    error
	Dur      time.Duration
}

var once = &sync.Once{}

//var httpResChan chan chan HttpProxyJobResponse
var httpProxyJobPool = sync.Pool{
	New: func() interface{} {
		return &httpProxyJob{}
	},
}

func initHttpProxyResChan(chanSize int) {
	once.Do(func() {
		log.Success("Res channel size:%d", chanSize)
		//httpResChan = make(chan chan HttpProxyJobResponse,chanSize)
		//for i:=0;i<chanSize;i++{
		//	httpResChan <- make(chan HttpProxyJobResponse)
		//}
	})
}
func NewHttpProxyJob(transport *http.Transport, q *http.Request, id string) *httpProxyJob {
	//if httpResChan == nil{
	//	initHttpProxyResChan(rcsize)
	//}
	job := httpProxyJobPool.Get().(*httpProxyJob)
	job.q = q
	job.m = id
	job.t = transport
	job.r = make(chan HttpProxyJobResponse)
	//job := &httpProxyJob{q:q,m:id,t:transport,r: make(chan HttpProxyJobResponse)}
	//job.setResChan()
	return job
}

type httpProxyJob struct {
	q *http.Request
	r chan HttpProxyJobResponse
	m string
	t *http.Transport
}

func (job *httpProxyJob) Id() string {
	return job.m
}
func (job *httpProxyJob) OnWorkerFull(dispatcher *NonBlockingDispatcher) {
	job.r <- HttpProxyJobResponse{nil, errors.New("worker繁忙"), 0}
}
func (job *httpProxyJob) Do() {
	//now := time.Now()
	//log.Info("%s,%s",job.t.RoundTrip,job.q.Method)

	res, err := job.t.RoundTrip(job.q)
	//1.加协程可快速释放worker,worker不论转发是否成功就回列
	// 优点:处理速度快,保证worker快速回,保证可用worker数
	// 缺点:假如后端高并发处理能力不足,则会造成雪崩效应，后端已经处理不过来了，这边还是持续由worker发送请求
	//2.不加协程可由worker数量+请求处理速度,利用workder堵塞,可适当保护后端服务,提升流量爆发期间的tps(后续请求会马上返回worker已满),当然worker数量的配置得和处理能力相当
	// 优点:量化处理能力，如果后端处理得过来，则增加worker数量，否则减少
	// 缺点:高并发会因为超时处理过多,job.R消费等待过长会发生worker不够用的情况
	//目前选择:2,处理能力可通过调整worker size大小来决定
	//TODO : 更多测试
	//go func(res *http.Response,err error) {
	job.r <- HttpProxyJobResponse{res, err, 0}
	//log.Info("%s -> transport请求时间消耗 : %s",job.Id(),time.Since(now))
	//}()
}
func (job *httpProxyJob) GetResChan() chan HttpProxyJobResponse {
	return job.r
}
func (job *httpProxyJob) Release() {
	<-job.r
	//job.putResChan()
	//httpProxyJobPool.Put(job)
}
