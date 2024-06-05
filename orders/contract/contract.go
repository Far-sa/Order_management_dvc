package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}
type OrderRepository interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
	UpdateOrder(ctx context.Context, id string, order *pb.Order) error
}
