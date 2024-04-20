package api

import (
	"fmt"
	"net/http"

	"github.com/olafstar/salejobs-api/internal/s3"
)

type APIServer struct {
	listenAddr string
	store Store
	requestQueueManager *RequestQueueManager
	cache *allCache
	s3 *s3.S3Client
}

func NewAPIServer (listenAddr string, db Store, rqm *RequestQueueManager, cache *allCache, s3 *s3.S3Client) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: db,
		requestQueueManager: rqm,
		cache: cache,
		s3: s3,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	s.SetupUserAPI(mux)
	s.SetupAdvertismentAPI(mux)
	s.SetupUtilsRoutes(mux)

	fmt.Printf("Server listening on http://localhost%s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}