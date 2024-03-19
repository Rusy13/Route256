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
	fmt.Println("Внутренняя ошибка сервера:")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pvzRepo := &repository.Pvz{
		PvzName: unm.PvzName,
		Address: unm.Address,
		Email:   unm.Email,
	}
	id, err := s.Repo.Add(req.Context(), pvzRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	fmt.Println("Внутренняя ошибка сервера:")
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("yutu")
	pvz, err := s.Repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(err)

	pvzJson, _ := json.Marshal(pvz)
	w.Write(pvzJson)
}

func (s *Server1) Update(w http.ResponseWriter, req *http.Request) {
	// Получаем ID из пути запроса
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Декодирование JSON из тела запроса в структуру addPvzRequest
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Удаление статьи из базы данных
	if err := s.Repo.Delete(req.Context(), keyInt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа об успешном удалении
	pvzJson := []byte("SUCCESS")
	w.Write(pvzJson)
}
