package adapter

import (
	"context"
	"testing"
	"xyz_multifinance/src/internal/creditlimit/domain/creditlimit"
	"xyz_multifinance/src/internal/shared/mysql"
	"xyz_multifinance/src/internal/testutils"

	"gopkg.in/go-playground/assert.v1"
)

func newCreditLimit() *creditlimit.CreditLimit {
	return &creditlimit.CreditLimit{
		CustomerID: "customer-uuid-1",
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

func TestRepository(t *testing.T) {
	t.Parallel()

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		t.Error("error create mysql connection")
	}

	testutils.SeedCustomer(t, db, "customer-uuid-1", "John doe", "23482342430")
	r := NewMysqlCreditLimitRepository(db)

	t.Run("TestCreate", func(t *testing.T) {
		t.Parallel()
		testCreate(t, r)
	})
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
				return newCreditLimit()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			customer := tc.Create(t)

			err := repository.Create(ctx, customer)
			assert.Equal(t, err, nil)
		})
	}

}
