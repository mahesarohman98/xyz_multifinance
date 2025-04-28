package query

import (
	"context"
	"xyz_multifinance/src/internal/customer/domain/customer"
	"xyz_multifinance/src/internal/shared/decorator"

	"github.com/sirupsen/logrus"
)

type GetCustomerByID struct {
	CustomerByID string
}

type GetCustomerByIDHandler decorator.QueryHandler[GetCustomerByID, *customer.Customer]

type GetCustomerByIDReadModel interface {
	GetCustomer(
		ctx context.Context,
		customerID string,
	) (*customer.Customer, error)
}

type getCustomerByIDHandler struct {
	readModel GetCustomerByIDReadModel
}

func NewGetCustomerByIDHandler(
	readModel GetCustomerByIDReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetCustomerByIDHandler {
	return decorator.ApplyQueryDecorators(
		getCustomerByIDHandler{
			readModel: readModel,
		},
		logger,
		metricsClient,
	)
}
func (h getCustomerByIDHandler) Handle(ctx context.Context, query GetCustomerByID) (*customer.Customer, error) {
	return h.readModel.GetCustomer(ctx, query.CustomerByID)
}
