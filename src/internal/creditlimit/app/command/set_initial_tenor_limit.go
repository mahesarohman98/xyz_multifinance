package command

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/decorator"

	"github.com/sirupsen/logrus"
)

type SetInitialTenorLimit struct {
	CustomerID string
	Tenors     []struct {
		MonthRange  int
		LimitAmount float64
	}
}

type SetInitialTenorLimitHandler decorator.CommandHandler[SetInitialTenorLimit]

type setInitialTenorLimitHandler struct {
	fc   creditlimit.Factory
	repo creditlimit.Repository
}

func NewSetInitialTenorLimitHandler(
	fc creditlimit.Factory,
	repo creditlimit.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) SetInitialTenorLimitHandler {
	return decorator.ApplyCommandDecorators(
		setInitialTenorLimitHandler{
			fc:   fc,
			repo: repo,
		},
		logger,
		metricsClient,
	)
}

func (h setInitialTenorLimitHandler) Handle(ctx context.Context, cmd SetInitialTenorLimit) error {
	creditLimitUser := h.fc.MustNewCreditLimit(cmd.CustomerID)
	for _, t := range cmd.Tenors {
		if err := creditLimitUser.AddTenor(t.MonthRange, t.LimitAmount); err != nil {
			return err
		}
	}

	return h.repo.Create(ctx, creditLimitUser)
}
