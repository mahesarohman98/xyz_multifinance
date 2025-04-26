package adapter

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	mysqlutil "xyz_multifinance/src/internal/shared/mysql"
)

type mysqlTenor struct {
	CustomerID string  `db:"customer_id"`
	Month      int     `db:"month"`
	Amount     float64 `db:"amount"`
}

func parseCreditLimit(creditLimit *creditlimit.CreditLimit) []mysqlTenor {
	tenors := []mysqlTenor{}
	for _, t := range creditLimit.Tenors {
		tenors = append(tenors, mysqlTenor{
			CustomerID: creditLimit.CustomerID,
			Month:      t.MonthRange,
			Amount:     t.LimitAmount,
		})
	}
	return tenors
}

type MysqlCreditLimitRepository struct {
	db *sqlx.DB
}

func NewMysqlCreditLimitRepository(db *sqlx.DB) *MysqlCreditLimitRepository {
	return &MysqlCreditLimitRepository{
		db: db,
	}
}

func (m MysqlCreditLimitRepository) Create(
	ctx context.Context,
	creditLimit *creditlimit.CreditLimit,
) error {
	for {
		err := m.create(ctx, creditLimit)
		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && val.Number == mysqlutil.MySQLDeadlockErrorCode {
			continue
		}
		return err
	}
}

func (m MysqlCreditLimitRepository) create(
	ctx context.Context,
	creditLimit *creditlimit.CreditLimit,
) error {
	mysqlTenors := parseCreditLimit(creditLimit)
	if len(mysqlTenors) <= 0 {
		return nil
	}

	query := `
		INSERT INTO 
			TenorLimits (customer_id, month, amount) 
		VALUES
	`
	args := []interface{}{}
	valuePlaceholders := ""

	for i, t := range mysqlTenors {
		if i > 0 {
			valuePlaceholders += ", "
		}
		valuePlaceholders += "(?, ?, ?)"
		args = append(args, t.CustomerID, t.Month, t.Amount)
	}

	fullQuery := query + valuePlaceholders

	_, err := m.db.ExecContext(ctx, fullQuery, args...)
	if err != nil {
		return errors.Wrap(err, "failed to batch insert TenorLimits")
	}

	return nil
}
