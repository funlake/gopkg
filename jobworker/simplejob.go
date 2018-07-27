package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
)

type simpleJob struct{

}
func (sj *simpleJob) Do(){
	//log.Info("helloworld")
	//time.Sleep(2)
}
func (sj *simpleJob)Id() string{
	return ""
}
func (sj *simpleJob)OnWorkerFull(){
	//log.Warning("Worker is full:%d",fullSize)
}
