package command

import "context"

type CreditLimitService interface {
	GetTotalUsedByCustomerAndTenor(
		ctx context.Context,
		customerID string,
		tenor int,
		forUpdate bool,
	) (float64, error)

	DecreaseLimit(
		ctx context.Context,
		customerID string,
		tenor int,
		totalBorowed float64,
	) error
}
