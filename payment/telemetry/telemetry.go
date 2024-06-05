package telemetry

import (
	"context"
	"fmt"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/payment/contract"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next contract.PaymentsService
}

func NewTelemetryMiddleware(next contract.PaymentsService) contract.PaymentsService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreatePayment: %v", o))

	return s.next.CreatePayment(ctx, o)
}
