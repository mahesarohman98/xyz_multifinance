package customer

import (
	"context"
)

type Repository interface {
	CreateOrUpdate(
		ctx context.Context,
		customer *Customer,
	) error
}
