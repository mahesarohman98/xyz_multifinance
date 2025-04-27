package command

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/decorator"

	"github.com/sirupsen/logrus"
)

type DecreaseLimit struct {
	CustomerID   string
	MonthRange   int
	TotalBorowed float64
}

type DecreaseLimitHandler decorator.CommandHandler[DecreaseLimit]

type decreaseLimitHandler struct {
	repo creditlimit.Repository
}

func NewDecreaseLimitHandler(
	repo creditlimit.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) DecreaseLimitHandler {
	return decorator.ApplyCommandDecorators(
		decreaseLimitHandler{
			repo: repo,
		},
		logger,
		metricsClient,
	)
}
func (h decreaseLimitHandler) Handle(ctx context.Context, cmd DecreaseLimit) error {
	if err := h.repo.Update(
		ctx,
		cmd.CustomerID,
		func(cl *creditlimit.CreditLimit) (*creditlimit.CreditLimit, error) {
			if err := cl.DecreaseLimit(cmd.MonthRange, cmd.TotalBorowed); err != nil {
				return nil, err
			}

			return cl, nil

		}); err != nil {
		return err
	}
	return nil
}
