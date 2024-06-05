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

func (s *service) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {

	id, err := s.orderRepository.CreateOrder(ctx, in, items)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID:         id,
		CustomerID: in.CustomerID,
		Status:     "pending",
		Items:      items,
	}

	return o, nil

}

func (s *service) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.Order, error) {
	return s.orderRepository.Get(ctx, in.OrderID, in.CustomerID)
}

func (s *service) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	err := s.orderRepository.UpdateOrder(ctx, order.ID, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, in *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(in.Items) == 0 {
		return nil, common.ErrNoItems
	}
	mergedItems := mergeItemsQuantities(in.Items)

	log.Println("", mergedItems)

	//* validate with stock
	//! Temp :
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1PNz51RxQqzMVLiGKLqohCOK",
			ID:       i.ItemID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil

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
