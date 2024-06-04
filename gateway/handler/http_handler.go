package handler

import (
	"errors"
	"log"
	"net/http"

	common "github.com/Far-sa/commons"
	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/gateway/gateway"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func New(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {

	//* static serving
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.handleGetOrder)
}

func (h *handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	orderID := r.PathValue("orderID")

	order, _ := h.gateway.GetOrder(r.Context(), orderID, customerID)

	//!! convert error
	// sErr := status.Convert(err)
	// if sErr != nil {
	// 	if sErr.Code() != codes.InvalidArgument {
	// 		common.WriteError(w, http.StatusBadRequest, sErr.Message())
	// 		return
	// 	}

	// }
	common.WriteJson(w, http.StatusOK, order)

}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("creating order")

	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJson(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	//TODO use validator library
	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	//!! convert error
	// sErr := status.Convert(err)
	// if sErr != nil {
	// 	if sErr.Code() != codes.InvalidArgument {
	// 		common.WriteError(w, http.StatusBadRequest, sErr.Message())
	// 		return
	// 	}

	// }

	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJson(w, http.StatusOK, o)

}

// * Helper function
func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items must not be empty")
	}

	for _, item := range items {
		if item.ItemID == "" {
			return errors.New("item ID is required")
		}
		if item.Quantity <= 0 {
			return errors.New("item must have a valid quantity")
		}
	}
	return nil
}

// ! handle with Echo
// func (h *handler) RegisterRoutes(e *echo.Echo) {
// 	e.POST("/api/customers/:customerID/orders", h.HandleCreateOrder)
// }

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
