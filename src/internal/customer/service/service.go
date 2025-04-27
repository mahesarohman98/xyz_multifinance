package service

import (
	"context"
	"xyz_multifinance/src/internal/customer/adapter"
	"xyz_multifinance/src/internal/customer/app"
	"xyz_multifinance/src/internal/customer/app/command"
	"xyz_multifinance/src/internal/customer/app/query"
	"xyz_multifinance/src/internal/customer/domain/customer"
	"xyz_multifinance/src/internal/shared/metrics"
	"xyz_multifinance/src/internal/shared/mysql"

	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := mysql.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	factory := customer.NewFactory(customer.FactoryConfig{
		WageLimit:  3500000,
		MinimumAge: 20,
	})
	repo := adapter.NewMysqlCustomerRepository(db)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			RegisterNewCustomer: command.NewRegisterNewCustomerHandler(factory, repo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerByID: query.NewGetCustomerByIDHandler(repo, logger, metricsClient),
		},
	}

}
