package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/utils"
)
type Dispatcher struct {
	workerPool chan chan WorkerJob
	jobQueue chan WorkerJob
}
func NewDispather(maxWorker int,queueSize int) *Dispatcher {
	dispatcher := &Dispatcher{}
	//流水线长度
	dispatcher.jobQueue = make(chan WorkerJob,queueSize)
	//流水线工人数量
	dispatcher.workerPool = make(chan chan WorkerJob,maxWorker)
	dispatcher.Run(maxWorker)
	return dispatcher
}
func (d *Dispatcher) Put(job WorkerJob) bool{
	select {
		case d.jobQueue <- job:
			return true
		default :
			return false
			//log.Error("job队列已满")
	}
	return false
}
func (d *Dispatcher) Run(maxWorker int){
	for i:=0 ;i<maxWorker;i++{
		worker := NewWorker(d.workerPool)
		worker.Ready()
	}
	utils.WrapGo(func() {
		d.Ready()
	})
}

//需做worker繁忙超时处理机制,否则会因为worker过于繁忙，
//症状是任务快速加入，然而worker繁忙导致管道堵塞，并发时，如果worker数量远小于并发，
//如不做超时处理,则后续请求会在<-d.workerPool处等待很长时间，而真正进入转发阶段的消耗时间却并不多
//所以后端服务无压力，而请求却堵在网关处!
func (d *Dispatcher) Ready(){
	for{
		select{
			case job := <-d.jobQueue:
				select {
					case jobChan := <-d.workerPool :
						jobChan <- job
					default:
						job.OnWorkerFull(d)
				}
		}
	}
}