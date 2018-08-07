package jobworker

import (
	"time"
	"github.com/funlake/gopkg/utils"
)

// blocking dispather
type BlockingDispatcher struct {
	workerPool chan chan WorkerNonBlockingJob
	jobQueue chan WorkerNonBlockingJob
}

func NewBlockingDispather(maxWorker int,queueSize int) *BlockingDispatcher {
	dispatcher := &BlockingDispatcher{}
	//流水线长度
	dispatcher.jobQueue = make(chan WorkerNonBlockingJob,queueSize)
	//流水线工人数量
	dispatcher.workerPool = make(chan chan WorkerNonBlockingJob,maxWorker)
	dispatcher.Run(maxWorker)
	//稍微等下worker启动
	time.Sleep(time.Nanosecond * 10)
	return dispatcher
}

func (d *BlockingDispatcher) Put(job WorkerNonBlockingJob) bool{
	d.jobQueue <- job
	return true
}
func (d *BlockingDispatcher) Run(maxWorker int){
	for i:=0 ;i<maxWorker;i++{
		worker := NewWorker(d.workerPool)
		worker.Ready()
	}
	utils.WrapGo(func() {
		d.Ready()
	},"dispather run")
}

//后端持续处理需要堵塞并后续执行，故无需超时处理
func (d *BlockingDispatcher) Ready(){
	for{
		select{
		case job := <-d.jobQueue:
			select {
			case jobChan := <-d.workerPool :
				jobChan <- job
			}
		}
	}
}