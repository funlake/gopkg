package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
	//"github.com/funlake/gopkg/utils/log"
)
func NewSimpleJob() *simpleJob {
	job := &simpleJob{}
	return job
}
type simpleJob struct{

}
func (sj *simpleJob) Do(){
	//log.Info("helloworld")
	//time.Sleep(2)
}
func (sj *simpleJob)Id() string{
	return ""
}
func (sj *simpleJob)OnWorkerFull(dispatcher *BlockingDispatcher){
	//log.Warning("Worker is full")
	//dispatcher.Put(sj)
}
