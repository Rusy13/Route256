package main

import (
	"HW1/api"
	"HW1/pkg/db"
	"HW1/pkg/repository/postgresql"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pvzRepo := postgresql.NewArticles(database)
	implementation := api.Server1{Repo: pvzRepo}

	http.Handle("/", api.CreateRouter(implementation))
	if err := http.ListenAndServe(api.Port, nil); err != nil {
		log.Fatal(err)
	}
}
