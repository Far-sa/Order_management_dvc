package repository

import (
	"context"
	"errors"

	pb "github.com/Far-sa/commons/api"
)

var orders = make([]*pb.Order, 0)

type repository struct {
	// MongoDB
}

func New() *repository {
	return &repository{}
}

func (r *repository) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest, items []*pb.Item) (string, error) {

	id := "24"
	orders = append(orders, &pb.Order{
		ID:         "24",
		CustomerID: in.CustomerID,
		Status:     "pending",
		Items:      items,
	})
	return id, nil
}

func (r *repository) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {

	for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	}

	return nil, errors.New("order not found")

}

func (r *repository) UpdateOrder(ctx context.Context, id string, order *pb.Order) error {
	for i, o := range orders {
		if o.ID == id {
			orders[i] = order
			return nil
		}
	}

	return errors.New("order not found")
}
