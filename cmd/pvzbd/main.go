package main

import (
	"HW1/api"
	config "HW1/internal/config"
	"HW1/internal/storage/db"
	"HW1/internal/storage/repository/postgresql"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	securePort   = ":9000"
	insecurePort = ":9001"
)

func main() {
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	log.Fatal("Failed to get current directory:", err)
	//}

	// Формируем путь к файлу .env с использованием текущей директории
	//envPath := filepath.Join(currentDir, "..", "..", ".env")
	//envPath := "/home/ubunto/Desktop/Route256/Route256DZ/HW1/.env" ////////////////

	//if err := godotenv.Load(envPath); err != nil {
	//	log.Fatal("No .env file found:", err)
	//}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to convert PORT to integer:", err)
	}
	config := config.StorageConfig{
		Host:     os.Getenv("HOST"),
		Port:     port,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("PASSWORD"),
		Database: os.Getenv("DBNAME"),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx, config)
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
	if err := http.ListenAndServeTLS(securePort, "api/server.crt", "api/server.key", secureMux); err != nil {
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
