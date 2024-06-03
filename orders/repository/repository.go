package repository

import "context"

type repository struct {
	// MongoDB
}

func New() *repository {
	return &repository{}
}

func (r *repository) CreateOrder(ctx context.Context) error {
	return nil
}
