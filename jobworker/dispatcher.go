package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
	"sync"
)
type Dispatcher interface {
	Put(job WorkerJob) bool
	Run(maxWork int)
	Ready()
	Metrics() sync.Map
}