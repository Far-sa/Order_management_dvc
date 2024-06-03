package contract

import "context"

type OrderService interface {
	CreateOrder(context.Context) error
}
type OrderRepository interface {
	CreateOrder(context.Context) error
}
