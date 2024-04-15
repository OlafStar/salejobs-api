package main

import (
	"github.com/olafstar/salejobs-api/internal/api"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db := api.InitDatabase()
	server := api.NewAPIServer(":4200", db)
	server.Run()
}