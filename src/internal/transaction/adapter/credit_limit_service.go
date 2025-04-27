package adapter

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/app"
	"xyz_multifinance/src/internal/creditlimit/app/command"
	"xyz_multifinance/src/internal/creditlimit/app/query"
)

type CreditLimitService struct {
	app app.Application
}

func NewCreditLimitService(app app.Application) CreditLimitService {
	return CreditLimitService{app: app}
}

func (s *CreditLimitService) GetTotalUsedByCustomerAndTenor(
	ctx context.Context,
	customerID string,
	tenor int,
	forUpdate bool,
) (float64, error) {
	return s.app.Queries.GetTotalUsedByCustomerAndTenor.Handle(ctx, query.GetTotalUsedByCustomerAndTenor{
		CustomerID: customerID,
		Tenor:      tenor,
		ForUpdate:  forUpdate,
	})
}

func (s *CreditLimitService) DecreaseLimit(
	ctx context.Context,
	customerID string,
	tenor int,
	totalBorowed float64,
) error {
	return s.app.Commands.DecreaseLimit.Handle(ctx, command.DecreaseLimit{
		CustomerID:   customerID,
		MonthRange:   tenor,
		TotalBorowed: totalBorowed,
	})
}
