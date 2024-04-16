package main

import (
	"github.com/olafstar/salejobs-api/internal/api"
	"github.com/stripe/stripe-go/v72"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	stripe.Key = "sk_test_51O9XshLphFr90yTInpBihTpHwMGcJ8WE2ZaJzKY2UNG2MFNGexMNNUAOYOhvI5ab6a82IusDYAI3w1y0TAEhPCt300zuqzMVaD"
	db := api.InitDatabase()
	server := api.NewAPIServer(":4200", db)
	server.Run()
}