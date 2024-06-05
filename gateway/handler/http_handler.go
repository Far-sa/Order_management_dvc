package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	common "github.com/Far-sa/commons"
	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/gateway/gateway"
	"github.com/Far-sa/gateway/param"
	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	order, err := h.gateway.GetOrder(ctx, orderID, customerID)

	rStatus := status.Convert(err)
	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())

		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

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

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	//TODO use validator library
	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())

		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := &param.CreateOrderResponse{
		Order:         o,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s/"+o.CustomerID, o.ID),
	}

	common.WriteJson(w, http.StatusOK, res)

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
