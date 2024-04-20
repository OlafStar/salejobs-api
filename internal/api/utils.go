package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/olafstar/salejobs-api/internal/middleware"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

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

func ReadFileFromRequest(r *http.Request) ([]byte, string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, "", &HTTPError{StatusCode: http.StatusBadRequest, Message: "Error parsing multipart form"}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, "", &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid file"}
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, "", &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Error reading file"}
	}

	return fileData, header.Filename, nil
}