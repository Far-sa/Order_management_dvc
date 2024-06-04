package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
}
type OrderRepository interface {
	CreateOrder(context.Context) error
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
}
