package command

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/decorator"
)

type DecreaseLimit struct {
	customerID   string
	monthRange   int
	totalBorowed float64
}

type DecreaseLimitHandler decorator.CommandHandler[DecreaseLimit]

type decreaseLimitHandler struct {
	repo creditlimit.Repository
}

func (h decreaseLimitHandler) Handle(ctx context.Context, cmd DecreaseLimit) error {
	if err := h.repo.Update(
		ctx,
		cmd.customerID,
		func(cl *creditlimit.CreditLimit) (*creditlimit.CreditLimit, error) {
			if err := cl.DecreaseLimit(cmd.monthRange, cmd.totalBorowed); err != nil {
				return nil, err
			}

			return cl, nil

		}); err != nil {
		return err
	}
	return nil
}
