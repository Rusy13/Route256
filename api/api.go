package postgresql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"Homework/internal/storage/repository"
)

const Port = ":9000"
const queryParamKey = "key"
const cacheTTL = time.Second * 15

type Server1 struct {
	Repo        repository.PvzRepo
	RedisClient *redis.Client
}

type addPvzRequest struct {
	PvzName string `json:"pvzname"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type addPvzResponse struct {
	ID      int64  `json:"id"`
	PvzName string `json:"pvzname"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func CreateRouter(implemetation Server1) *mux.Router {
	router := mux.NewRouter()
	router.Use(BasicAuth, LoggingMiddleware)
	router.HandleFunc("/pvz", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.CreatePvz(w, req)
		default:
			fmt.Println("error")
		}
	})
	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetPVZByID(w, req)
		case http.MethodDelete:
			implemetation.DeletePvz(w, req)
		case http.MethodPut:
			implemetation.UpdatePvz(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

func (s *Server1) CreatePvz(w http.ResponseWriter, req *http.Request) {
	req.BasicAuth()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	pvzRepo := &repository.Pvz{
		PvzName: unm.PvzName,
		Address: unm.Address,
		Email:   unm.Email,
	}
	id, err := s.Repo.Add(req.Context(), pvzRepo)
	if err != nil {
		http.Error(w, "Failed to add pvz", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &addPvzResponse{
		ID:      id,
		PvzName: pvzRepo.PvzName,
		Address: pvzRepo.Address,
		Email:   pvzRepo.Email,
	}
	pvzJson, _ := json.Marshal(resp)
	w.Write(pvzJson)
}

func (s *Server1) GetPVZByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "Invalid request parameter", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	pvzJson, status := s.get(req.Context(), keyInt)

	w.WriteHeader(status)
	w.Write(pvzJson)
}

func validateGetByID(key int64) bool {
	if key <= 0 {
		return false
	}
	return true
}

func (s *Server1) get(ctx context.Context, key int64) ([]byte, int) {
	// Чистка redis (FlushAll)
	cacheKey := fmt.Sprintf("pvz:%d", key)
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, http.StatusInternalServerError
		}
	} else {
		log.Println("Данные из кэша")
		return cachedData, http.StatusOK
	}

	article, err := s.Repo.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Printf("Article with key %s not found in database", key)
			return nil, http.StatusNotFound
		}

		log.Printf("Error retrieving data from database for key %s: %v", key, err)
		return nil, http.StatusInternalServerError
	}

	log.Println("Successfully retrieved data from database")

	pvzJson, err := json.Marshal(article)
	if err != nil {
		log.Printf("Error serializing article data for key %s: %v", key, err)
		return nil, http.StatusInternalServerError
	}

	err = s.RedisClient.Set(ctx, cacheKey, pvzJson, cacheTTL).Err()
	if err != nil {
		log.Printf("Error caching data to Redis for key %s: %v", key, err)
		return nil, http.StatusInternalServerError
	}

	log.Println("Data successfully cached in Redis")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pvzJson, http.StatusOK
}

func (s *Server1) UpdatePvz(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "Invalid request parameter", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		//fmt.Println(err)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatedPvz := &repository.Pvz{
		ID:      int64(keyInt),
		PvzName: unm.PvzName,
		Address: unm.Address,
		Email:   unm.Email,
	}

	if err := s.Repo.Update(req.Context(), keyInt, updatedPvz); err != nil {
		http.Error(w, "Failed to update pvz", http.StatusInternalServerError)
		return
	}

	pvzJson, _ := json.Marshal(updatedPvz)
	w.Write(pvzJson)
}

func (s *Server1) DeletePvz(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "Invalid request parameter", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := s.Repo.Delete(req.Context(), keyInt); err != nil {
		http.Error(w, "Failed to delete pvz", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}
