package throttler

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	var tests = []struct {
		Desc       string
		Jobs       []string
		MaxWorkers int
		TotalJobs  int
	}{
		{
			"Standard implementation",
			[]string{"job01", "job02", "job03", "job04", "job05", "job06", "job07", "job08", "job09", "job10",
				"job11", "job12", "job13", "job14", "job15", "job16", "job17", "job18", "job19", "job20",
				"job21", "job22", "job23", "job24", "job25", "job26", "job27", "job28", "job29", "job30",
				"job31", "job32", "job33", "job34", "job35", "job36", "job37", "job38", "job39", "job40",
				"job41", "job42", "job43", "job44", "job45", "job46", "job47", "job48", "job49", "job50"},
			5,
			-1,
		}, {
			"Incorrectly has 0 as TotalWorkers",
			[]string{"job01", "job02", "job03", "job04", "job05", "job06", "job07", "job08", "job09", "job10",
				"job11", "job12", "job13", "job14", "job15", "job16", "job17", "job18", "job19", "job20",
				"job21", "job22", "job23", "job24", "job25", "job26", "job27", "job28", "job29", "job30",
				"job31", "job32", "job33", "job34", "job35", "job36", "job37", "job38", "job39", "job40",
				"job41", "job42", "job43", "job44", "job45", "job46", "job47", "job48", "job49", "job50"},
			5,
			0,
		}, {
			"More workers than jobs",
			[]string{"job01", "job02", "job03", "job04", "job05", "job06", "job07", "job08", "job09", "job10",
				"job11", "job12", "job13", "job14", "job15", "job16", "job17", "job18", "job19", "job20",
				"job21", "job22", "job23", "job24", "job25", "job26", "job27", "job28", "job29", "job30",
				"job31", "job32", "job33", "job34", "job35", "job36", "job37", "job38", "job39", "job40",
				"job41", "job42", "job43", "job44", "job45", "job46", "job47", "job48", "job49", "job50"},
			50000,
			-1,
		},
	}

	for _, test := range tests {
		totalJobs := len(test.Jobs)
		if test.TotalJobs != -1 {
			totalJobs = test.TotalJobs
		}
		th := New(test.MaxWorkers, totalJobs)
		for _, job := range test.Jobs {
			go func(job string, th *Throttler) {
				defer th.Done(nil)
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			}(job, th)
			th.Throttle()
		}
	}
}

func TestThrottleWithErrors(t *testing.T) {
	var tests = []struct {
		Desc       string
		Jobs       []string
		MaxWorkers int
		TotalJobs  int
	}{
		{
			"Standard implementation",
			[]string{"job01", "job02", "job03", "job04", "job05", "job06", "job07", "job08", "job09", "job10",
				"job11", "job12", "job13", "job14", "job15", "job16", "job17", "job18", "job19", "job20",
				"job21", "job22", "job23", "job24", "job25", "job26", "job27", "job28", "job29", "job30",
				"job31", "job32", "job33", "job34", "job35", "job36", "job37", "job38", "job39", "job40",
				"job41", "job42", "job43", "job44", "job45", "job46", "job47", "job48", "job49", "job50"},
			5,
			-1,
		},
	}

	for _, test := range tests {
		totalJobs := len(test.Jobs)
		if test.TotalJobs != -1 {
			totalJobs = test.TotalJobs
		}
		th := New(test.MaxWorkers, totalJobs)
		for _, job := range test.Jobs {
			go func(job string, th *Throttler) {
				jobNum, _ := strconv.ParseInt(job[len(job)-2:], 10, 8)
				var err error
				if jobNum%2 != 0 {
					err = fmt.Errorf("Error on %s", job)
				}
				defer th.Done(err)

				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			}(job, th)
			th.Throttle()
		}
		if len(th.Err()) != totalJobs/2 {
			t.Fatal()
		}
	}
}

func TestThrottlePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal()
		}
	}()
	New(0, 100)
}
