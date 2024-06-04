package service

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	panic("not implemented")
}
