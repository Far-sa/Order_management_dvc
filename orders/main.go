package main

import (
	"context"

	"github.com/Far-sa/order/repository"
	"github.com/Far-sa/order/service"
)

func main() {
	repo := repository.New()
	svc := service.New(repo)

	svc.CreateOrder(context.Background())
}
