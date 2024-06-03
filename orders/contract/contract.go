package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type OrderService interface {
	CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.Order, error)
}
type OrderRepository interface {
	CreateOrder(context.Context) error
}
