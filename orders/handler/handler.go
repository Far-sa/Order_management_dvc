package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/order/contract"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service contract.OrderService
	ch      *amqp.Channel //* may pass the events in service layer
}

// !
func NewGRPC(grpcServer *grpc.Server, service contract.OrderService, ch *amqp.Channel) {
	handler := &grpcHandler{service: service, ch: ch}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, in)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order received! order %v:", in)

	//TODO move to service
	q, err := h.ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpan.End()

	//TODO implement validation as separate adapter or service
	items, err := h.service.ValidateOrder(amqpContext, in)
	if err != nil {
		return nil, err
	}

	order, err := h.service.CreateOrder(amqpContext, in, items)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	headers := broker.InjectAMQPHeaders(amqpContext)

	h.ch.PublishWithContext(amqpContext, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	return order, nil
}

func (h *grpcHandler) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, order)
}
