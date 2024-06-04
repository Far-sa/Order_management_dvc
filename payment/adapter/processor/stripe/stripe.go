package stripeProcessor

import (
	"context"
	"fmt"

	common "github.com/Far-sa/commons"
	pb "github.com/Far-sa/commons/api"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var gatewayHTTPAddr = common.EnvString("GATEWAY_HTTP_ADDRESS", "http://localhost:8080")

type stripeProcessor struct {
}

func NewProcessor() *stripeProcessor {
	return &stripeProcessor{}
}

func (s *stripeProcessor) CreatePaymentLink(ctx context.Context, order *pb.Order) (string, error) {
	fmt.Printf("creating payment link for order %+v\n", order)

	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, order.CustomerID, order.ID)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", gatewayHTTPAddr)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PricID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", nil
	}

	return result.URL, nil
}
