package customer

import (
	"context"
)

type Repository interface {
	Create(
		ctx context.Context,
		customer *Customer,
	) error
}
