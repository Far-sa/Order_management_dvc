package service

import (
	"context"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/payment/contract"
)

type service struct {
	processor contract.PaymentProcessor
}

func NewService(processor contract.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(ctx, o)
	if err != nil {
		return "", err
	}

	//TODO update order with link
	return link, nil
}
