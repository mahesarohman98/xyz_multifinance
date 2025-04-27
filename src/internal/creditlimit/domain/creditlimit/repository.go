package creditlimit

import "context"

type Repository interface {
	GetCreditLimit(
		ctx context.Context,
		customerID string,
	) (*CreditLimit, error)

	Create(
		ctx context.Context,
		creditLimit *CreditLimit,
	) error

	Update(
		ctx context.Context,
		customerID string,
		fn func(*CreditLimit) (*CreditLimit, error),
	) error
}
