package api

import (
	"net/http"

	"github.com/olafstar/salejobs-api/internal/env"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

// Error implements error.
func (a apiFunc) Error() string {
	panic("unimplemented")
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type ApiError struct {
	Error string `json:"error"`
}

func (s *APIServer) SetupUserAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", s.makeHTTPHandleFunc(s.HandleLogin))
	mux.HandleFunc("/api/auth/register", s.makeHTTPHandleFunc(s.HandleRegister))
	mux.HandleFunc("/api/iam", s.makeHTTPHandleFunc(ProtectedRequest(s.Iam)))
}

func (s *APIServer) SetupAdvertismentAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/advertisements", s.makeHTTPHandleFunc(s.handleAdvertisements))
	mux.HandleFunc("/api/advertisements/counter", s.makeHTTPHandleFunc(s.handleAdvertisementsCounter))
	mux.HandleFunc("/api/advertisements/{id}", s.makeHTTPHandleFunc(s.handleSpecificAdvertisement))
}

//TODO: Delete this before prod
func (s *APIServer) SetupUtilsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/test", s.makeHTTPHandleFunc(func(w http.ResponseWriter, r *http.Request) error {
			if r.Method == "POST" {
					fileData, fileName, err := ReadFileFromRequest(r)

					if err != nil {
						return err
					}

					if err := s.s3.UploadData(fileData, env.GoEnv("R2_BUCKET_NAME"), fileName); err != nil {
							return err
					}

					return WriteJSON(w, http.StatusOK, "File uploaded successfully")
			}

			return &HTTPError{StatusCode: http.StatusMethodNotAllowed, Message: "Method not allowed"}
	}))
}