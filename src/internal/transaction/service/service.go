package service

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/service"
	"xyz_multifinance/src/internal/shared/metrics"
	"xyz_multifinance/src/internal/shared/mysql"
	"xyz_multifinance/src/internal/transaction/adapter"
	"xyz_multifinance/src/internal/transaction/app"
	"xyz_multifinance/src/internal/transaction/app/command"

	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	creditLimitApp := service.NewApplication(ctx)
	creditLimitService := adapter.NewCreditLimitService(creditLimitApp)

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	repo := adapter.NewMysqlTransactionRepository(db)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			SubmitLoad: command.NewSubmitLoanHandler(repo, &creditLimitService, logger, metricsClient),
		},
		Queries: app.Queries{},
	}

}
