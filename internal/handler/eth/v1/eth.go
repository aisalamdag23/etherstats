package v1

import (
	"encoding/json"
	"net/http"

	"github.com/aisalamdag23/etherstats/internal/domain"
	"github.com/aisalamdag23/etherstats/internal/handler"
	"github.com/gorilla/mux"
)

type server struct {
	service domain.Service
}

type apiProblem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func NewServer(service domain.Service) handler.Handler {
	return &server{
		service: service,
	}
}

func (s *server) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/eth/{id}", handler.Restrict(http.MethodGet, s.GetEth))
}

func (s *server) GetEth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	resp, err := s.service.Get(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(apiProblem{
			Title:  http.StatusText(http.StatusInternalServerError),
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
