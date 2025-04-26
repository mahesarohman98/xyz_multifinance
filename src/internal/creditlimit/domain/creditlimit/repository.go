package creditlimit

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		creditLimit *CreditLimit,
	) error
}
