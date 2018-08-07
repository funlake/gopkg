package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
)
type Dispatcher interface {
	Put(job WorkerNonBlockingJob) bool
	Run(maxWork int)
	Ready()
}
