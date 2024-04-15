package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/olafstar/salejobs-api/internal/middleware"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: []string{"https://example.com", "*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}
	return middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
			if err := f(w, r); err != nil {
					var httpErr *HTTPError
					if errors.As(err, &httpErr) {
							WriteJSON(w, httpErr.StatusCode, ApiError{Error: httpErr.Message})
					} else {
							WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
					}
					return
			}
	}, middleware.Logging(), middleware.CORS(corsConfig))
}