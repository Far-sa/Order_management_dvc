package service

import (
	"context"
	"log"

	common "github.com/Far-sa/commons"
	pb "github.com/Far-sa/commons/api"
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

func (s *service) ValidateOrder(ctx context.Context, in *pb.CreateOrderRequest) error {
	if len(in.Items) == 0 {
		return common.ErrNoItems
	}
	mergedItems := mergeItemsQuantities(in.Items)

	log.Println("", mergedItems)

	//* validate with repository
	return nil

}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ItemID == item.ItemID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
