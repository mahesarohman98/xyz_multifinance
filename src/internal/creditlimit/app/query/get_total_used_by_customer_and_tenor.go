package query

import (
	"context"
	"xyz_multifinance/src/internal/shared/decorator"

	"github.com/sirupsen/logrus"
)

type GetTotalUsedByCustomerAndTenor struct {
	CustomerID string
	Tenor      int
	ForUpdate  bool
}

type GetTotalUsedByCustomerAndTenorHandler decorator.QueryHandler[GetTotalUsedByCustomerAndTenor, float64]

type GetTotalUsedByCustomerReadModel interface {
	GetTotalUsedByCustomerAndTenor(ctx context.Context, customerID string, tenor int, forUpdate bool) (float64, error)
}

type getTotalUsedByCustomerAndTenorHandler struct {
	readModel GetTotalUsedByCustomerReadModel
}

func NewGetTotalUsedByCustomerAndTenorHandler(
	readModel GetTotalUsedByCustomerReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetTotalUsedByCustomerAndTenorHandler {
	return decorator.ApplyQueryDecorators(
		getTotalUsedByCustomerAndTenorHandler{
			readModel: readModel,
		},
		logger,
		metricsClient,
	)
}

func (h getTotalUsedByCustomerAndTenorHandler) Handle(ctx context.Context, query GetTotalUsedByCustomerAndTenor) (float64, error) {
	return h.readModel.GetTotalUsedByCustomerAndTenor(ctx, query.CustomerID, query.Tenor, true)
}
