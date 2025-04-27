package service

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/adapter"
	"xyz_multifinance/src/internal/creditlimit/app"
	"xyz_multifinance/src/internal/creditlimit/app/command"
	"xyz_multifinance/src/internal/creditlimit/app/query"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/metrics"
	"xyz_multifinance/src/internal/shared/mysql"

	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := mysql.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	factory := creditlimit.NewFactory()
	repo := adapter.NewMysqlCreditLimitRepository(db)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			SetInitialTenorLimit: command.NewSetInitialTenorLimitHandler(factory, repo, logger, metricsClient),
			DecreaseLimit:        command.NewDecreaseLimitHandler(repo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetTotalUsedByCustomerAndTenor: query.NewGetTotalUsedByCustomerAndTenorHandler(repo, logger, metricsClient),
			GetCreditLimitByCustomerID:     query.NewGetCreditLimitByCustomerHandler(repo, logger, metricsClient),
		},
	}
}
