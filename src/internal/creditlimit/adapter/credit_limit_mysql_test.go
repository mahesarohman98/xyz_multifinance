package adapter

import (
	"context"
	"testing"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/mysql"
	"xyz_multifinance/src/internal/testutils"

	"gopkg.in/go-playground/assert.v1"
)

func newCreditLimit(customerID string) *creditlimit.CreditLimit {
	return &creditlimit.CreditLimit{
		CustomerID: customerID,
		Tenors: []creditlimit.Tenor{
			{
				MonthRange:  1,
				LimitAmount: 100,
			},
			{
				MonthRange:  2,
				LimitAmount: 200,
			},
			{
				MonthRange:  3,
				LimitAmount: 300,
			},
			{
				MonthRange:  6,
				LimitAmount: 600,
			},
		},
	}
}

func testCreate(t *testing.T, repository creditlimit.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name   string
		Create func(*testing.T) *creditlimit.CreditLimit
	}{
		{
			Name: "create customer",
			Create: func(t *testing.T) *creditlimit.CreditLimit {
				return newCreditLimit("customer-uuid-1")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			creditLimit := tc.Create(t)

			err := repository.Create(ctx, creditLimit)
			assert.Equal(t, err, nil)
		})
	}

}

func assertCreditLimitInRepository(ctx context.Context, t *testing.T, repo creditlimit.Repository, creditLimit *creditlimit.CreditLimit) {
	assert.NotEqual(t, creditLimit, nil)

	creditLimitFromRepo, err := repo.GetCreditLimit(ctx, creditLimit.CustomerID)
	assert.Equal(t, err, nil)

	assert.Equal(t, creditLimit, creditLimitFromRepo)
}

func testUpdate(t *testing.T, repository creditlimit.Repository) {
	t.Helper()
	ctx := context.Background()

	testCreditLimit := newCreditLimit("customer-uuid-2")

	err := repository.Create(ctx, testCreditLimit)
	assert.Equal(t, err, nil)

	var expectedCreditLimit *creditlimit.CreditLimit
	err = repository.Update(
		ctx,
		"customer-uuid-2",
		func(creditLimit *creditlimit.CreditLimit) (*creditlimit.CreditLimit, error) {
			if err := creditLimit.DecreaseLimit(3, 100); err != nil {
				return nil, err
			}
			expectedCreditLimit = creditLimit
			return creditLimit, nil
		})
	assert.Equal(t, err, nil)

	assertCreditLimitInRepository(ctx, t, repository, expectedCreditLimit)
}

func TestRepository(t *testing.T) {
	t.Parallel()

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		t.Error("error create mysql connection")
	}

	testutils.SeedCustomer(t, db, "customer-uuid-1", "John doe", "23482342430")
	testutils.SeedCustomer(t, db, "customer-uuid-2", "Jane doe", "2398334r130")
	r := NewMysqlCreditLimitRepository(db)

	t.Run("TestCreate", func(t *testing.T) {
		t.Parallel()
		testCreate(t, r)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		t.Parallel()
		testUpdate(t, r)
	})
}
