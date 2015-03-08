// Throttler fills the gap between sync.WaitGroup and manually monitoring your goroutines
// with channels. The API is almost identical to Wait Groups, but it allows you to set
// a max number of workers that can be running simultaneously. It uses channels internally
// to block until a job completes by calling Done() or until all jobs have been completed.
//
// Compare the Throttler example to the sync.WaitGroup example http://golang.org/pkg/sync/#example_WaitGroup
//
// See a fully functional example on the playground at http://bit.ly/throttler-docs
package throttler

type Throttler struct {
	maxWorkers    int
	workerCount   int
	totalJobs     int
	jobsStarted   int
	jobsCompleted int
	doneChan      chan struct{}
}

// New returns a Throttler that will govern the max number of workers and will
// work with the total number of jobs. It panics if maxWorkers < 1.
func New(maxWorkers, totalJobs int) *Throttler {
	if maxWorkers < 1 {
		panic("maxWorkers has to be at least 1")
	}
	return &Throttler{
		maxWorkers: maxWorkers,
		totalJobs:  totalJobs,
		doneChan:   make(chan struct{}, totalJobs),
	}
}

// Throttle works similarly to sync.WaitGroup, except inside your goroutine dispatch
// loop rather than after. It will not block until the number of active workers
// matches the max number of workers designated in the call to NewThrottler or
// all of the jobs have been dispatched. It stops blocking when Done has been called
// as many times as totalJobs.
func (t *Throttler) Throttle() {
	if t.totalJobs < 1 {
		return
	}
	t.jobsStarted++
	t.workerCount++

	if t.workerCount == t.maxWorkers {
		<-t.doneChan
		t.jobsCompleted++
		t.workerCount--
	}

	if t.jobsStarted == t.totalJobs {
		for t.jobsCompleted < t.totalJobs {
			<-t.doneChan
			t.jobsCompleted++
		}
	}
}

// Done lets Throttler know that a job has been completed so that another worker
// can be activated. If Done is called less times than totalJobs,
// Throttle will block forever
func (t *Throttler) Done() {
	t.doneChan <- struct{}{}
}
