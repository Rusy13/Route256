package main

import (
	postgresql "Homework/api_grpc"
	pb "Homework/protos/gen/go/app"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	api "Homework/api"
	config "Homework/internal/config"
	"Homework/internal/storage/db"
	pp "Homework/internal/storage/repository/postgresql"
)

const (
	securePort   = ":9000"
	insecurePort = ":9001"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	portStr := os.Getenv("POSTGRES_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Error converting port to integer")
	}
	config := config.StorageConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DBNAME"),
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")

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

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		panic("Error converting REDIS_DB to int")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	pvzRepo := pp.NewPvzRepo(database)
	implementation := api.Server1{Repo: pvzRepo,
		RedisClient: rdb}

	// сервера GRPC
	//repo := postgresql.NewPvzRepository() // Замените этот вызов на вашу реализацию репозитория
	go func() {
		//repo := postgresql.NewPvzRepository() // Замените этот вызов на вашу реализацию репозитория
		server := &postgresql.Server{
			Repo: pvzRepo,
		}

		// Создание TCP соединения на порту 10000
		lis, err := net.Listen("tcp", ":10000")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Создание GRPC сервера
		grpcServer := grpc.NewServer()

		// Регистрация сервера GRPC
		pb.RegisterPvzServiceServer(grpcServer, server)

		// Запуск GRPC сервера
		log.Println("Starting gRPC server on port :10000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	go serveSecure(implementation)
	go metric()
	serveInsecure()

}

func metric() {
	log.Printf("metrics on 2112")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func serveSecure(implementation api.Server1) {
	secureMux := http.NewServeMux()
	secureMux.Handle("/", api.CreateRouter(implementation))

	log.Printf("Listening on port %s...\n", securePort)
	if err := http.ListenAndServeTLS(securePort, "../../api/server.crt", "../../api/server.key", secureMux); err != nil {
		log.Fatal(err)
	}
}

func serveInsecure() {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostParts := strings.Split(req.Host, ":")
		host := hostParts[0]

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
