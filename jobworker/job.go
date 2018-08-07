package jobworker
type WorkerNonBlockingJob interface{
	Do()
	Id() string
	OnWorkerFull(dispatcher *NonBlockingDispatcher)
}
type WorkerBlockingJob interface{
	Do()
	Id() string
	OnWorkerFull(dispatcher *BlockingDispatcher)
}

