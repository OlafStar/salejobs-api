package api

import (
	"fmt"
	"net/http"

)

type APIServer struct {
	listenAddr string
	store Store
	requestQueueManager *RequestQueueManager
}

func NewAPIServer (listenAddr string, db Store, rqm *RequestQueueManager) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: db,
		requestQueueManager: rqm,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	s.SetupUserAPI(mux)
	s.SetupAdvertismentAPI(mux)
	

	fmt.Printf("Server listening on http://localhost%s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}