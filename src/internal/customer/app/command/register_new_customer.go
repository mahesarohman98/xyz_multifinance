package command

import (
	"context"
	"time"
	"xyz_multifinance/src/internal/customer/domain/customer"
	"xyz_multifinance/src/internal/shared/decorator"
	"xyz_multifinance/src/internal/shared/errors"

	"github.com/sirupsen/logrus"
)

// RegisterNewCustomer
// Todo photo upload
type RegisterNewCustomer struct {
	ID           string
	NIK          string
	Fullname     string
	LegalName    string
	PlaceOfBirth string
	DateOfBirth  time.Time
	Wage         float64
	PhotoURL     string
	KTPURL       string
	Today        time.Time
}

type RegisterNewCustomerHandler decorator.CommandHandler[RegisterNewCustomer]

type registerNewCustomerHandler struct {
	fc   customer.Factory
	repo customer.Repository
}

func NewRegisterNewCustomerHandler(
	fc customer.Factory,
	repo customer.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RegisterNewCustomerHandler {
	return decorator.ApplyCommandDecorators(
		registerNewCustomerHandler{
			fc:   fc,
			repo: repo,
		},
		logger,
		metricsClient,
	)
}

func (h registerNewCustomerHandler) Handle(ctx context.Context, cmd RegisterNewCustomer) error {
	customer, err := h.fc.RegisterNewCustomer(
		cmd.ID,
		cmd.NIK,
		cmd.Fullname,
		cmd.LegalName,
		cmd.PlaceOfBirth,
		cmd.DateOfBirth,
		cmd.Wage,
		cmd.PhotoURL,
		cmd.KTPURL,
		cmd.Today,
	)
	if err != nil {
		return errors.NewIncorrectInputError(err.Error(), "incorrect-register-customer")
	}

	if err := h.repo.Create(ctx, customer); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-create-customer")
	}

	return nil
}
