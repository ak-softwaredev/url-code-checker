package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var noMux error = fmt.Errorf("no mux provided")

type API struct{}

type jsonConfig struct {
	Links []string            `json:"links"`
	Codes map[string]Response `json:"responses"`
}

func SetRoutes(mux *http.ServeMux) error {
	if mux == nil {
		return noMux
	}

	api := API{}

	mux.HandleFunc("POST /check", api.postCheck)
	mux.HandleFunc("GET /check/{id}", api.getCheck)
	mux.HandleFunc("GET /health", api.health)

	return nil
}

func (a *API) postCheck(w http.ResponseWriter, r *http.Request) {
	links := jsonConfig{}

	if err := json.NewDecoder(r.Body).Decode(&links); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	response, err := Check(r.Context(), links.Links)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	links.Codes = response
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
