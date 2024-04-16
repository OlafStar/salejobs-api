package api

import (
	"net/http"
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
	mux.HandleFunc("/api/advertisments", s.makeHTTPHandleFunc(s.handleAdvertisments))
}
