package handler

import (
	"context"
	"log"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/order/contract"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	orderService contract.OrderService
	pb.UnimplementedOrderServiceServer
}

func NewGRPC(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New Order received")
	order := &pb.Order{
		ID: "24",
	}
	return order, nil
}
