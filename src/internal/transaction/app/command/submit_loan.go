package command

import (
	"context"
	"errors"
	"xyz_multifinance/src/internal/shared/decorator"
	transaction "xyz_multifinance/src/internal/transaction/domain/trasaction"

	"github.com/sirupsen/logrus"
)

type Loan struct {
	ID             string
	ContractNumber string
	OTR            float64
	AmountInterest float64
	AssetName      string
}

type SubmitLoan struct {
	Source struct {
		ID         string
		ExternalID string
	}
	CustomerID string
	Tenor      int
	Loans      []Loan
}

type SubmitLoanHandler decorator.CommandHandler[SubmitLoan]

type submitLoadHandler struct {
	repo                    transaction.Repository
	creditLimitQueryService CreditLimitService
}

func NewSubmitLoanHandler(
	repo transaction.Repository,
	CreditLimitService CreditLimitService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) SubmitLoanHandler {
	return decorator.ApplyCommandDecorators(
		submitLoadHandler{
			repo:                    repo,
			creditLimitQueryService: CreditLimitService,
		},
		logger,
		metricsClient,
	)
}

func (h submitLoadHandler) Handle(ctx context.Context, cmd SubmitLoan) error {
	if err := h.repo.Create(
		ctx,
		func() ([]transaction.Transaction, error) {
			totalTenorUsed, err := h.creditLimitQueryService.GetTotalUsedByCustomerAndTenor(
				ctx, cmd.CustomerID, cmd.Tenor, true,
			)
			if err != nil {
				return []transaction.Transaction{}, err
			}

			transactions := []transaction.Transaction{}
			totalBorrowAmount := float64(0)

			for _, loan := range cmd.Loans {
				tx, err := transaction.NewTransaction(
					loan.ID,
					loan.ContractNumber,
					cmd.CustomerID,
					cmd.Source,
					cmd.Tenor,
					loan.OTR,
					loan.AmountInterest,
					loan.AssetName,
				)
				if err != nil {
					return []transaction.Transaction{}, err
				}

				totalBorrowAmount += tx.TotalBorowed()
				transactions = append(transactions, *tx)
			}

			totalTenorUsed -= totalBorrowAmount
			if totalTenorUsed <= 0 {
				return []transaction.Transaction{}, errors.New("total tenor used exceed")
			}

			if err := h.creditLimitQueryService.DecreaseLimit(
				ctx, cmd.CustomerID, cmd.Tenor, totalBorrowAmount,
			); err != nil {
				return []transaction.Transaction{}, err
			}

			return transactions, nil
		}); err != nil {
		return err
	}

	return nil
}
