package main

import (
	"HW1/api"
	"HW1/pkg/db"
	"HW1/pkg/repository/postgresql"
	"context"
	"log"
	"net/http"
)

const (
	securePort   = ":9000"
	insecurePort = ":9001"
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

	go serveSecure(implementation)
	serveInsecure(implementation)
}

func serveSecure(implementation api.Server1) {
	secureMux := http.NewServeMux()
	secureMux.Handle("/", api.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", securePort)
	if err := http.ListenAndServeTLS(securePort, "/home/ubunto/Desktop/Route256/Route256DZ/HW1/api/server.crt", "/home/ubunto/Desktop/Route256/Route256DZ/HW1/api/server.key", secureMux); err != nil {
		log.Fatal(err)
	}
}

func serveInsecure(implementation api.Server1) {
	insecureMux := http.NewServeMux()
	insecureMux.Handle("/", api.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", insecurePort)
	if err := http.ListenAndServe(insecurePort, insecureMux); err != nil {
		log.Fatal(err)
	}
}
