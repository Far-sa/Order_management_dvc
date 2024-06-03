package handler

import (
	"log"
	"net/http"
)

type handler struct {
	// gateway
}

func New() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("creating order")
}
