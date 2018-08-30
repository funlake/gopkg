package jobworker

import (
	"time"
	"github.com/funlake/gopkg/utils"
)

// blocking dispather
type BlockingDispatcher struct {
	workerPool chan chan WorkerJob
	jobQueue chan WorkerJob
	//failQueue chan WorkerJob
}

func NewBlockingDispather(maxWorker int,queueSize int) *BlockingDispatcher {
	dispatcher := &BlockingDispatcher{}
	//流水线长度
	dispatcher.jobQueue = make(chan WorkerJob,queueSize)
	//流水线工人数量
	dispatcher.workerPool = make(chan chan WorkerJob,maxWorker)
	//失败队列
	//dispatcher.failQueue = make(chan WorkerJob,queueSize * 2)
	dispatcher.Run(maxWorker)
	//稍微等下worker启动
	time.Sleep(time.Nanosecond * 10)
	return dispatcher
}

func (d *BlockingDispatcher) Put(job WorkerJob) bool{
	select {
		case d.jobQueue <- job:
			return true
		default :
			//d.failQueue <- job
			return false
	}

	return true
}
func (d *BlockingDispatcher) Run(maxWorker int){
	for i:=0 ;i<maxWorker;i++{
		worker := NewWorker(d.workerPool)
		worker.Ready()
	}
	utils.WrapGo(func() {
		//go d.Recover()
		d.Ready()
	},"dispather run")
}
//func (d *BlockingDispatcher) Recover(){
//	for {
//		select {
//			case job := <- d.failQueue:
//				d.jobQueue <- job
//			default :
//		}
//		time.Sleep(time.Second * 1)
//	}
//}
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