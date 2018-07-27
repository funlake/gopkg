package jobworker
type WorkerJob interface{
	Do()
	Id() string
	OnWorkerFull()
}

