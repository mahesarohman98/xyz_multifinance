package customer

import (
	"context"
)

type Repository interface {
	GetCustomer(
		ctx context.Context,
		customerID string,
	) (*Customer, error)

	Create(
		ctx context.Context,
		customer *Customer,
	) error
}
