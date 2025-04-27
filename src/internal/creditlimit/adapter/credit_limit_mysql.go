package adapter

import (
	"context"
	"database/sql"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	mysqlutil "xyz_multifinance/src/internal/shared/mysql"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlTenor struct {
	CustomerID string  `db:"customer_id"`
	Month      int     `db:"month"`
	Amount     float64 `db:"amount"`
	UsedAmount float64 `db:"used_amount"`
}

func parseCreditLimit(creditLimit *creditlimit.CreditLimit) []mysqlTenor {
	tenors := []mysqlTenor{}
	for _, t := range creditLimit.Tenors {
		tenors = append(tenors, mysqlTenor{
			CustomerID: creditLimit.CustomerID,
			Month:      t.MonthRange,
			Amount:     t.LimitAmount,
			UsedAmount: t.UsedAmount,
		})
	}
	return tenors
}

func UnmarshallToDomain(records []mysqlTenor) *creditlimit.CreditLimit {
	if len(records) <= 0 {
		return nil
	}
	creditLimit := &creditlimit.CreditLimit{
		CustomerID: records[0].CustomerID,
		Tenors:     make([]creditlimit.Tenor, 0, len(records)),
	}

	for _, record := range records {
		tenor := creditlimit.Tenor{
			MonthRange:  record.Month,
			LimitAmount: record.Amount,
			UsedAmount:  record.UsedAmount,
		}
		creditLimit.Tenors = append(creditLimit.Tenors, tenor)
	}

	return creditLimit
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
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		mysqlutil.FinishTransaction(err, tx)
	}()

	if err := m.upsertCreditLimit(ctx, tx, creditLimit); err != nil {
		return err
	}

	return nil
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
type sqlContextGetter interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func (m MysqlCreditLimitRepository) GetTotalUsedByCustomerAndTenor(
	ctx context.Context,
	customerID string,
	tenor int,
	forUpdate bool,
) (float64, error) {
	query := `SELECT used_amount FROM TenorLimits tl WHERE tl.customer_id = ? AND tl.month = ?`
	if forUpdate {
		query += " FOR UPDATE"
	}

	rows, err := m.db.QueryContext(ctx, query, customerID, tenor)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalUsed float64
	for rows.Next() {
		if err := rows.Scan(&totalUsed); err != nil {
			return 0, err
		}
	}

	return totalUsed, nil
}

func (m MysqlCreditLimitRepository) GetCreditLimit(
	ctx context.Context,
	customerID string,
) (*creditlimit.CreditLimit, error) {
	return m.getCreditLimit(ctx, m.db, customerID, false)
}

func (m MysqlCreditLimitRepository) getCreditLimit(
	ctx context.Context,
	db sqlContextGetter,
	customerID string,
	forUpdate bool,
) (*creditlimit.CreditLimit, error) {
	query := "SELECT * FROM TenorLimits tl WHERE tl.customer_id = ?"
	if forUpdate {
		query += " FOR UPDATE"
	}

	rows, err := db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenors []mysqlTenor
	for rows.Next() {
		var t mysqlTenor
		if err := rows.Scan(&t.CustomerID, &t.Month, &t.Amount, &t.UsedAmount); err != nil {
			return nil, err
		}
		tenors = append(tenors, t)
	}

	if len(tenors) <= 0 {
		return nil, sql.ErrNoRows
	}

	creditLimit := UnmarshallToDomain(tenors)

	return creditLimit, nil

}

func (m MysqlCreditLimitRepository) upsertCreditLimit(ctx context.Context, tx *sql.Tx, creditLimitToUpdate *creditlimit.CreditLimit) error {
	query := `
		INSERT INTO 
			TenorLimits (customer_id, month, amount, used_amount) 
		VALUES
	`
	args := []interface{}{}
	valuePlaceholders := ""

	updatedDbCreditLimit := parseCreditLimit(creditLimitToUpdate)
	if len(updatedDbCreditLimit) <= 0 {
		return nil
	}
	for i, t := range updatedDbCreditLimit {
		if i > 0 {
			valuePlaceholders += ", "
		}
		valuePlaceholders += "(?, ?, ?, ?)"
		args = append(args, t.CustomerID, t.Month, t.Amount, t.UsedAmount)
	}

	fullQuery := query + valuePlaceholders + `
		ON DUPLICATE KEY UPDATE
			used_amount = VALUES(used_amount)
	`

	_, err := tx.ExecContext(ctx, fullQuery, args...)
	if err != nil {
		return errors.Wrap(err, "unable to upsert credit limit")
	}

	return nil
}

func (m MysqlCreditLimitRepository) Update(
	ctx context.Context,
	customerID string,
	fn func(*creditlimit.CreditLimit) (*creditlimit.CreditLimit, error),
) error {
	for {
		err := m.update(ctx, customerID, fn)
		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && val.Number == mysqlutil.MySQLDeadlockErrorCode {
			continue
		}
		return err
	}
}

func (m MysqlCreditLimitRepository) update(
	ctx context.Context,
	customerID string,
	fn func(*creditlimit.CreditLimit) (*creditlimit.CreditLimit, error),
) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		mysqlutil.FinishTransaction(err, tx)
	}()

	existingCreditLimit, err := m.getCreditLimit(ctx, tx, customerID, true)
	if err != nil {
		return err
	}

	updatedCreditLimit, err := fn(existingCreditLimit)
	if err != nil {
		return err
	}

	if updatedCreditLimit == nil {
		return errors.New("updatedCreditLimit must not nil")
	}

	if err := m.upsertCreditLimit(ctx, tx, updatedCreditLimit); err != nil {
		return err
	}

	return nil
}
