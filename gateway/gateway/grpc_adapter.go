package gateway

import (
	"context"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/commons/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
}
