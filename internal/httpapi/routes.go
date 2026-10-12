package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type API struct{}

type jsonConfig struct {
	Links     []string `json:"links"`
	Codes     []Result `json:"responses"`
	TotalTime int64    `json:"total_ms"`
}

func SetRoutes(mux *http.ServeMux) {
	api := API{}

	mux.HandleFunc("POST /api/v1/check", api.postCheck)
	mux.HandleFunc("GET /api/v1/check/{id}", api.getCheck)
	mux.HandleFunc("GET /health", api.health)
}

func (a *API) postCheck(w http.ResponseWriter, r *http.Request) {
	links := jsonConfig{}

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024) // 16 KB
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&links); err != nil {
		var maxBytes *http.MaxBytesError
		if errors.As(err, &maxBytes) {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	response, total_ms := Check(r.Context(), links.Links)
	links.Codes = response
	links.TotalTime = total_ms

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(links); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (a *API) getCheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
