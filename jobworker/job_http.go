package jobworker

import (
	"time"
	"net/http"
	"errors"
	"github.com/funlake/gopkg/utils/log"
	"sync"
)
type HttpProxyJobResponse struct{
	Response *http.Response
	Error error
	Dur time.Duration
}
var httpProxyJobPool = sync.Pool{
	New: func() interface{}{
		return &httpProxyJob{}
	},
}
func NewHttpProxyJob(transport *http.Transport,q *http.Request,r chan HttpProxyJobResponse,id string) *httpProxyJob {
	job := httpProxyJobPool.Get().(*httpProxyJob)
	job.q = q
	job.m = id
	job.r = r
	job.t = transport
	return job
	//return &httpProxyJob{q:q,r:r,m:id,t:transport}
}

type httpProxyJob struct {
	q *http.Request
	r chan HttpProxyJobResponse
	m string
	t *http.Transport
}
func (job *httpProxyJob) Id() string{
	return job.m
}
func (job *httpProxyJob) OnWorkerFull(dispatcher *Dispatcher){
	job.r <- HttpProxyJobResponse{nil, errors.New("worker繁忙"),0}
}
func(job *httpProxyJob) Do() {
	now := time.Now()
	res,err := job.t.RoundTrip(job.q)
	//1.加协程可快速释放worker,worker不论转发是否成功就回列
	// 优点:处理速度快,保证worker快速回,保证可用worker数
	// 缺点:假如后端高并发处理能力不足,则会造成雪崩效应，后端已经处理不过来了，这边还是持续由worker发送请求
	//2.不加协程可由worker数量+请求处理速度,利用workder堵塞,可适当保护后端服务,提升流量爆发期间的tps(后续请求会马上返回worker已满),当然worker数量的配置得和处理能力相当
	// 优点:量化处理能力，如果后端处理得过来，则增加worker数量，否则减少
	// 缺点:高并发会因为超时处理过多,job.R消费等待过长会发生worker不够用的情况
	//目前选择:2,处理能力可通过调整worker size大小来决定
	//TODO : 更多测试
	//go func(res *http.Response,err error) {
	job.r <- HttpProxyJobResponse{res, err, time.Since(now)}
	log.Info("%s -> 请求时间消耗 : %s",job.Id(),time.Since(now))
	//}()
}
func (job *httpProxyJob) GetResponse() chan HttpProxyJobResponse {
	r := job.r
	return r
}
func (job *httpProxyJob) ResetPool()  {
	httpProxyJobPool.Put(job)
}