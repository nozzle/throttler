package throttler

type httpPkg struct{}

func (httpPkg) Get(url string) {}

var http httpPkg

// This example fetches several URLs concurrently,
// using a WaitGroup to block until all the fetches are complete.
func ExampleWaitGroup() {
}

// This example fetches several URLs concurrently,
// using a Throttler to block until all the fetches are complete.
// Compare to http://golang.org/pkg/sync/#example_WaitGroup
func ExampleThrottler() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	// Create a new Throttler that will get 2 urls at a time
	t := NewThrottler(2, len(urls))
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Let Throttler know when the goroutine completes
			// so it can dispatch another worker
			defer t.Done()
			// Fetch the URL.
			http.Get(url)
		}(url)
		// Pauses until a worker is available or all jobs have been completed
		t.Throttle()
	}
}
