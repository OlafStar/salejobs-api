package api

import (
	"fmt"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store Store
	requestQueueManager *RequestQueueManager
	cache *allCache
}

func NewAPIServer (listenAddr string, db Store, rqm *RequestQueueManager, cache *allCache) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: db,
		requestQueueManager: rqm,
		cache: cache,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	s.SetupUserAPI(mux)
	s.SetupAdvertismentAPI(mux)


	fmt.Printf("Server listening on http://localhost%s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}