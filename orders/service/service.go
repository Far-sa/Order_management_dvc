package service

import (
	"context"

	"github.com/Far-sa/order/contract"
)

type service struct {
	orderRepository contract.OrderRepository
}

func New(orderRepository contract.OrderRepository) *service {
	return &service{orderRepository: orderRepository}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}
