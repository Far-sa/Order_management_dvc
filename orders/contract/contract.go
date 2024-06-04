package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}
type OrderRepository interface {
	CreateOrder(context.Context) error
}
