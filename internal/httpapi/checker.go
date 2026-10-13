package httpapi

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Job struct {
	Context context.Context
	URL     string
	Reply   chan<- Result
	Done    func()
}

type Result struct {
	URL  string `json:"url"`
	Code int    `json:"code"`
	Time int64  `json:"ms"`
}

var jobs = make(chan Job, 4*1024)
var client = http.Client{
	Timeout: 5 * time.Second,
}

func Check(ctx context.Context, links []string) <-chan Result {
	results := make(chan Result, len(links))

	wg := sync.WaitGroup{}

	for _, link := range links {
		valid := validateURL(link)
		if !valid {
			results <- Result{
				URL:  link,
				Code: -1,
				Time: 0,
			}
			continue
		}

		job := Job{
			Context: ctx,
			URL:     link,
			Reply:   results,
			Done:    func() { wg.Done() },
		}

		wg.Add(1)
		jobs <- job
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

var wereWorkersInitialized = false

func Init(workerAmount int) {
	if wereWorkersInitialized {
		return
	}

	for range workerAmount {
		go worker()
	}
}

func validateURL(link string) (valid bool) {
	parsed, err := url.ParseRequestURI(link)
	if err != nil ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Host == "" {
		return false
	}
	return true
}

func worker() {
	for job := range jobs {
		select {
		case <-job.Context.Done():
			job.Reply <- Result{
				URL:  job.URL,
				Code: -1,
			}
			job.Done()
			continue
		default:
			result := Result{}
			link := job.URL
			result.URL = link

			req, err := http.NewRequestWithContext(
				job.Context,
				http.MethodGet,
				link,
				nil,
			)
			if err != nil {
				result.Code = -1
				log.Printf("httpapi.Check -> %s: %s", link, err.Error())
				job.Reply <- result
				job.Done()
				continue
			}

			t0 := time.Now()
			r, err := client.Do(req)
			t1 := time.Now()
			result.Time = t1.Sub(t0).Milliseconds()
			if err != nil {
				result.Code = -1
				log.Printf("httpapi.Check -> %s: %s", link, err.Error())
				job.Reply <- result
				job.Done()
				continue
			}
			code := r.StatusCode

			result.Code = code

			r.Body.Close()
			job.Reply <- result

			job.Done()
		}
	}
}
