package main

import (
	"github.com/olafstar/salejobs-api/internal/api"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db := api.InitDatabase()
	rqm := api.NewRequestQueueManager(10, 10)
	defer rqm.Shutdown()
	server := api.NewAPIServer(":4200", db, rqm)
	server.Run()
}