package handler

import (
	"log"
	"net/http"

	common "github.com/Far-sa/commons"
	pb "github.com/Far-sa/commons/api"
)

type handler struct {
	client pb.OrderServiceClient
}

func New(client pb.OrderServiceClient) *handler {
	return &handler{client: client}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("creating order")

	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJson(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

}

//! handle with Echo
// func (h *handler) HandleCreateOrder(c echo.Context) error {
// 	log.Println("creating order")

// 	customerID := c.Param("customerID")
// 	var items []*pb.ItemsWithQuantity
// 	if err := c.Bind(&items); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	_, err := h.client.CreateOrder(c.Request().Context(), &pb.CreateOrderRequest{
// 		CustomerID: customerID,
// 		Items:      items,
// 	})
// 	if err != nil {
// 		// Assuming common.WriteError is a custom function that writes an error response.
// 		// You should replace this with appropriate error handling for your application.
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusCreated, map[string]string{"status": "order created"})
// }
