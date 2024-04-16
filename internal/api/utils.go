package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/olafstar/salejobs-api/internal/middleware"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// func (s *APIServer) makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
// 	corsConfig := middleware.CORSConfig{
// 		AllowedOrigins: []string{"https://example.com", "*"},
// 		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
// 		AllowedHeaders: []string{"Content-Type", "X-Requested-With"},
// 		AllowCredentials: true,
// 	}
// 	return middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
// 			if err := f(w, r); err != nil {
// 					var httpErr *HTTPError
// 					if errors.As(err, &httpErr) {
// 							WriteJSON(w, httpErr.StatusCode, ApiError{Error: httpErr.Message})
// 					} else {
// 							WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
// 					}
// 					return
// 			}
// 	}, middleware.Logging(), middleware.CORS(corsConfig))
// }

func (s *APIServer) makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: []string{"https://example.com", "*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}
	return middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		errc := make(chan error, 1)

		job := Job{
			Fn: func() error {
				err := f(w, r)
				return err
			},
			Errc: errc,
		}

		s.requestQueueManager.EnqueueJob(job)

		err := <-errc
		if err != nil {
			fmt.Println("Handling error from job")
			fmt.Println(err)
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				WriteJSON(w, httpErr.StatusCode, ApiError{Error: httpErr.Message})
			} else {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Internal server error"})
			}
		}
	}, middleware.Logging(), middleware.CORS(corsConfig))
}