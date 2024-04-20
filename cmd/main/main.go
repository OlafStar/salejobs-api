package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/olafstar/salejobs-api/internal/api"
	"github.com/olafstar/salejobs-api/internal/s3"
)

func main(){
	db := api.InitDatabase()
	s3 := s3.InitS3Client()
	rqm := api.NewRequestQueueManager(10, 10)
	c := api.NewCache()
	defer rqm.Shutdown()
	server := api.NewAPIServer(":4200", db, rqm, c, s3)
	server.Run()
}