package api

import (
	"HW1/pkg/repository"
	"HW1/pkg/repository/postgresql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

const Port = ":9000"
const queryParamKey = "key"

type Server1 struct {
	Repo *postgresql.PvzRepo
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
	router.Use(BasicAuth)
	router.HandleFunc("/pvz", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.Create(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetByID(w, req)
		case http.MethodDelete:
			implemetation.Delete(w, req)
		case http.MethodPut:
			implemetation.Update(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

func (s *Server1) Create(w http.ResponseWriter, req *http.Request) {
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

func (s *Server1) GetByID(w http.ResponseWriter, req *http.Request) {
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
	pvz, err := s.Repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, "Pvz not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	pvzJson, _ := json.Marshal(pvz)
	w.WriteHeader(http.StatusOK)
	w.Write(pvzJson)
}

func (s *Server1) Update(w http.ResponseWriter, req *http.Request) {
	// Получаем ID из пути запроса
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "Invalid request parameter", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Декодирование JSON из тела запроса в структуру addPvzRequest
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		//fmt.Println(err)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Создание объекта статьи с обновленными данными
	updatedPvz := &repository.Pvz{
		PvzName: unm.PvzName,
		Address: unm.Address,
		Email:   unm.Email,
	}

	// Обновление статьи в базе данных
	if err := s.Repo.Update(req.Context(), keyInt, updatedPvz); err != nil {
		http.Error(w, "Failed to update pvz", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа об успешном обновлении
	//w.WriteHeader(http.StatusOK)
	pvzJson, _ := json.Marshal(updatedPvz)
	w.Write(pvzJson)
}

func (s *Server1) Delete(w http.ResponseWriter, req *http.Request) {
	// Получаем ID из пути запроса
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		http.Error(w, "Invalid request parameter", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Удаление статьи из базы данных
	if err := s.Repo.Delete(req.Context(), keyInt); err != nil {
		http.Error(w, "Failed to delete pvz", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа об успешном удалении
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}
