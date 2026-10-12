package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Code int   `json:"code"`
	Time int64 `json:"ms"`
}

func testLink(ctx context.Context, url string) (Response, error) {
	select {
	case <-ctx.Done():
		return Response{}, ctx.Err()
	default:
		fmt.Println("Testing link:", url)
		response := Response{}

		client := http.Client{}
		t0 := time.Now()
		r, err := client.Get(url)
		if err != nil {
			return Response{}, err
		}
		t1 := time.Now()
		defer r.Body.Close()

		code := r.StatusCode

		response.Code = code
		response.Time = t1.Sub(t0).Milliseconds()

		return response, nil
	}
}

func Check(ctx context.Context, urls []string) (map[string]Response, error) {
	results := make(map[string]Response)

	for _, url := range urls {
		response, err := testLink(ctx, url)
		if err != nil {
			return map[string]Response{}, err
		}
		results[url] = response
	}

	return results, nil
}
