package transaction

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		fn func() ([]Transaction, error),
	) error
}
