package handler

import (
	"carinfo/internal/service"
	"net/http"
)

type Handler struct {
	CarService service.CarService
}

func NewHandler(s service.CarService) Handler {
	return Handler{
		CarService: s,
	}
}

func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/cars", h.handleListCars)
	mux.HandleFunc("/car/{id}", h.handleGetCarById)
	mux.HandleFunc("/car/new", h.handleCreateNewCar)
	mux.HandleFunc("/car/edit/{id}", h.handleEditCarById)
	mux.HandleFunc("/car/delete/{id}", h.handleDeleteCarById)
	return mux
}
