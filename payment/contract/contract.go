package contract

import (
	"context"

	pb "github.com/Far-sa/commons/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}

type PaymentProcessor interface {
	CreatePaymentLink(context.Context, *pb.Order) (string, error)
}