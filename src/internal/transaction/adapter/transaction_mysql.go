package adapter

import (
	"context"
	mysqlutil "xyz_multifinance/src/internal/shared/mysql"
	transaction "xyz_multifinance/src/internal/transaction/domain/trasaction"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlTransaction struct {
	TransactionID     string  `db:"transaction_id"`
	ContractNumber    string  `db:"contract_number"`
	CustomerID        string  `db:"customer_id"`
	ExternalID        string  `db:"external_id"`
	SourceID          string  `db:"source_id"`
	OTR               float64 `db:"otr"`
	AdminFee          float64 `db:"admin_fee"`
	TotalBorowed      float64 `db:"total_borowed"`
	InstallmentAmount float64 `db:"installment_amount"`
	AmountOfInterest  float64 `db:"amount_of_interest"`
	AssetName         string  `db:"asset_name"`
}

func parse(transactions []transaction.Transaction) []mysqlTransaction {
	mysqlTransactions := []mysqlTransaction{}
	for _, t := range transactions {
		mysqlTransactions = append(mysqlTransactions, mysqlTransaction{
			TransactionID:     t.ID,
			ContractNumber:    t.ContractNumber,
			CustomerID:        t.CustomerID,
			ExternalID:        t.Source.ExternalID,
			SourceID:          t.Source.ID,
			OTR:               t.OTR,
			AdminFee:          t.AdminFee(),
			TotalBorowed:      t.TotalBorowed(),
			InstallmentAmount: t.InstallmentAmount(),
			AmountOfInterest:  t.AmountInterest,
			AssetName:         t.AssetName,
		})
	}
	return mysqlTransactions

}

type MysqlTransactionRepository struct {
	db *sqlx.DB
}

func NewMysqlTransactionRepository(db *sqlx.DB) *MysqlTransactionRepository {
	return &MysqlTransactionRepository{
		db: db,
	}
}

func (m MysqlTransactionRepository) Create(
	ctx context.Context,
	fn func() ([]transaction.Transaction, error),
) error {
	for {
		err := m.create(ctx, fn)
		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && val.Number == mysqlutil.MySQLDeadlockErrorCode {
			continue
		}
		return err
	}
}

func (m MysqlTransactionRepository) create(
	ctx context.Context,
	fn func() ([]transaction.Transaction, error),
) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		mysqlutil.FinishTransaction(err, tx)
	}()

	transactions, err := fn()
	if err != nil {
		return err
	}

	mysqlTransactions := parse(transactions)
	if len(mysqlTransactions) <= 0 {
		return nil
	}

	query := `
		INSERT INTO 
			Transactions (
				transaction_id, contract_number, customer_id, external_id,
				source_id, otr, admin_fee, total_borowed, installment_amount,
				amount_of_interest, asset_name
			) 
		VALUES
	`
	args := []interface{}{}
	valuePlaceholders := ""

	for i, t := range mysqlTransactions {
		if i > 0 {
			valuePlaceholders += ", "
		}
		valuePlaceholders += `
			 (
				?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?
			) 
		`
		args = append(
			args, t.TransactionID, t.ContractNumber, t.CustomerID, t.ExternalID,
			t.SourceID, t.OTR, t.AdminFee, t.TotalBorowed, t.InstallmentAmount,
			t.AmountOfInterest, t.AssetName,
		)
	}

	fullQuery := query + valuePlaceholders

	_, err = tx.ExecContext(ctx, fullQuery, args...)
	if err != nil {
		return errors.Wrap(err, "failed to batch insert transactions")
	}

	return nil
}
