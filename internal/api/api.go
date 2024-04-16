package api

import (
	"fmt"
	"net/http"

)

type APIServer struct {
	listenAddr string
	store Store
}

func NewAPIServer (listenAddr string, db Store) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: db,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	s.SetupUserAPI(mux)
	s.SetupAdvertismentAPI(mux)
	SetupStripeRoutes(mux)

	fmt.Printf("Server listening on http://localhost%s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}