package query

import (
	"context"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/decorator"

	"github.com/sirupsen/logrus"
)

type GetCreditLimitByCustomerID struct {
	CustomerID string
}

type GetCreditLimitByCustomerIDHandler decorator.QueryHandler[GetCreditLimitByCustomerID, *creditlimit.CreditLimit]

type GetCreditLimitReadModel interface {
	GetCreditLimit(
		ctx context.Context,
		customerID string,
	) (*creditlimit.CreditLimit, error)
}

type getCreditLimitByCutsomerIDhandler struct {
	readModel GetCreditLimitReadModel
}

func NewGetCreditLimitByCustomerHandler(
	readModel GetCreditLimitReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetCreditLimitByCustomerIDHandler {
	return decorator.ApplyQueryDecorators(
		getCreditLimitByCutsomerIDhandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

func (h getCreditLimitByCutsomerIDhandler) Handle(ctx context.Context, query GetCreditLimitByCustomerID) (*creditlimit.CreditLimit, error) {
	return h.readModel.GetCreditLimit(ctx, query.CustomerID)
}
