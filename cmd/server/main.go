package main

import (
	"log"
	"net/http"

	"github.com/ak-softwaredev/url-code-checker/internal/httpapi"
)

func main() {
	httpMu := http.NewServeMux()

	httpapi.SetRoutes(httpMu)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpMu))
}
