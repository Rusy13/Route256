package main

import (
	api "HW1/api"
	config "HW1/internal/config"
	"HW1/internal/storage/db"
	pp "HW1/internal/storage/repository/postgresql"
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

	//port, err := strconv.Atoi(os.Getenv("PORT"))
	//if err != nil {
	//	log.Fatal("Failed to convert PORT to integer:", err)
	//}
	//config := config.StorageConfig{
	//	Host:     os.Getenv("HOST"),
	//	Port:     port,
	//	Username: os.Getenv("POSTGRES_USER"),
	//	Password: os.Getenv("PASSWORD"),
	//	Database: os.Getenv("DBNAME"),
	//}

	config := config.StorageConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "1111",
		Database: "Route",
	}
	//-------------------------------------------------------------------------------------------------------------
	var brokers = []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
	}
	//
	//api.ConsumerGroupExample(brokers)
	//go func() {
	//api.ConsumerGroupExample(brokers)
	//}()

	go func() {
		api.ConsumerExample(brokers)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pvzRepo := pp.NewPvzRepo(database)
	implementation := api.Server1{Repo: pvzRepo}

	go serveSecure(implementation)
	serveInsecure()
}

func serveSecure(implementation api.Server1) {
	secureMux := http.NewServeMux()
	secureMux.Handle("/", api.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", securePort)
	if err := http.ListenAndServeTLS(securePort, "../../api/server.crt", "../../api/server.key", secureMux); err != nil {
		//if err := http.ListenAndServeTLS(securePort, "api/server.crt", "api/server.key", secureMux); err != nil {
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
