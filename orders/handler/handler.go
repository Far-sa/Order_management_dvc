package handler

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/order/contract"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service contract.OrderService
	ch      *amqp.Channel
}

// !
func NewGRPC(grpcServer *grpc.Server, service contract.OrderService, ch *amqp.Channel) {
	handler := &grpcHandler{service: service, ch: ch}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order received! order %v:", in)

	order := &pb.Order{
		ID: "24",
	}

	q, err := h.ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	h.ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return order, nil
}
