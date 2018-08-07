package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/utils"
	"time"
)
type Dispatcher interface {
	Put(job WorkerJob) bool
	Run(maxWork int)
	Ready()
}
