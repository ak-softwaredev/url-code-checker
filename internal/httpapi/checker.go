package httpapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Result struct {
	URL  string `json:"url"`
	Code int    `json:"code"`
	Time int64  `json:"ms"`
}

var client = http.Client{
	Timeout: 10 * time.Second,
}

func Check(ctx context.Context, links []string) <-chan Result {
	jobs := make(chan string, len(links))
	results := make(chan Result, len(links))
	// results := make([]Result, 0, len(links))

	for _, link := range links {
		cutLink := validateAndCutURL(link)
		if cutLink == "" {
			results <- Result{
				URL:  link,
				Code: -1,
				Time: 0,
			}
			continue
		}

		jobs <- link
	}
	close(jobs)

	var wg = sync.WaitGroup{}

	for range 20 {
		wg.Add(1)
		go worker(ctx, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func validateAndCutURL(link string) (cutLink string) {
	parsed, err := url.ParseRequestURI(link)
	if err != nil ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Host == "" {
		return ""
	}
	if len(parsed.RequestURI()) == 1 {
		cutLink = parsed.Host + parsed.Path
	} else {
		cutLink = parsed.Host + parsed.RequestURI()
	}
	return cutLink
}

func worker(
	ctx context.Context,
	jobs <-chan string,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for link := range jobs {
		fmt.Println("Testing link:", link)
		result := Result{}
		result.URL = link

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
		if err != nil {
			result.Code = -1
			log.Printf("httpapi.Check -> %s: %s", link, err.Error())
			results <- result
			continue
		}

		t0 := time.Now()
		r, err := client.Do(req)
		t1 := time.Now()
		result.Time = t1.Sub(t0).Milliseconds()
		if err != nil {
			result.Code = -1
			log.Printf("httpapi.Check -> %s: %s", link, err.Error())
			results <- result
			continue
		}
		code := r.StatusCode

		result.Code = code

		r.Body.Close()

		results <- result
	}
}
