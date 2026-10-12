package httpapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Code int   `json:"code"`
	Time int64 `json:"ms"`
}

type Result struct {
	URL string `json:"url"`
	Response
}

func Check(ctx context.Context, links []string) ([]Result, int64) {
	results := make([]Result, 0, len(links))

	t0 := time.Now()
	for _, link := range links {
		cutLink := validateAndCutURL(link)
		if cutLink == "" {
			continue
		}

		response, err := testLink(ctx, link)
		if err != nil {
			log.Printf("httpapi.Check -> %s: %s", link, err.Error())
		}
		result := Result{
			URL:      cutLink,
			Response: response,
		}
		results = append(results, result)
	}
	t1 := time.Now()

	return results, t1.Sub(t0).Milliseconds()
}

func testLink(ctx context.Context, link string) (Response, error) {
	fmt.Println("Testing link:", link)
	response := Response{}

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)

	t0 := time.Now()
	r, err := client.Do(req)
	t1 := time.Now()
	response.Time = t1.Sub(t0).Milliseconds()
	if err != nil {
		response.Code = -1
		return response, err
	}
	defer r.Body.Close()

	code := r.StatusCode

	response.Code = code

	return response, nil

}

func validateAndCutURL(link string) (cutLink string) {
	parsed, err := url.ParseRequestURI(link)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return ""
	}
	if len(parsed.RequestURI()) == 1 {
		cutLink = parsed.Host + parsed.Path
	} else {
		cutLink = parsed.Host + parsed.RequestURI()
	}
	return cutLink
}
