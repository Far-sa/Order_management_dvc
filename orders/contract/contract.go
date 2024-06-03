package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type OrderService interface {
	CreateOrder(context.Context) error
	ValidateOrder(ctx context.Context, in *pb.CreateOrderRequest) error
}
type OrderRepository interface {
	CreateOrder(context.Context) error
}
