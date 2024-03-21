package main

import (
	"HW1/api"
	"HW1/pkg/db"
	"HW1/pkg/repository/postgresql"
	"context"
	"log"
	"net/http"
	"strings"
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
	serveInsecure()
}

func serveSecure(implementation api.Server1) {
	secureMux := http.NewServeMux()
	secureMux.Handle("/", api.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", securePort)
	if err := http.ListenAndServeTLS(securePort, "/home/ubunto/Desktop/Route256/Route256DZ/HW1/api/server.crt", "/home/ubunto/Desktop/Route256/Route256DZ/HW1/api/server.key", secureMux); err != nil {
		log.Fatal(err)
	}
}

func serveInsecure() {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostParts := strings.Split(req.Host, ":")
		host := hostParts[0]

		// Формируем целевой URL с портом 9000 для HTTPS
		target := "https://" + host + securePort + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, target, http.StatusTemporaryRedirect)
	})

	log.Printf("Listening on port %s...\n", insecurePort)
	if err := http.ListenAndServe(insecurePort, redirectHandler); err != nil {
		log.Fatal(err)
	}
}
