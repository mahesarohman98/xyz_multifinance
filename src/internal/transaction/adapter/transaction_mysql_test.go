package adapter

import (
	"context"
	"testing"
	"xyz_multifinance/src/internal/shared/mysql"
	"xyz_multifinance/src/internal/testutils"
	transaction "xyz_multifinance/src/internal/transaction/domain/trasaction"

	"gopkg.in/go-playground/assert.v1"
)

func newTransactions() []transaction.Transaction {
	return []transaction.Transaction{
		{
			ID:             "txid-1",
			ContractNumber: "contract-no-1",
			CustomerID:     "customer-tx-1",
			Source: transaction.Source{
				ID:         "sourceid-1",
				ExternalID: "1",
			},
			Tenor:          6,
			OTR:            8000000,
			AmountInterest: 1200,
			AssetName:      "sepeda",
		},
		{
			ID:             "txid-2",
			ContractNumber: "contract-no-2",
			CustomerID:     "customer-tx-1",
			Source: transaction.Source{
				ID:         "sourceid-1",
				ExternalID: "1",
			},
			Tenor:          6,
			OTR:            2000000,
			AmountInterest: 1200,
			AssetName:      "mesin cuci",
		},
	}
}

func testCreate(t *testing.T, repository transaction.Repository) {
	t.Helper()

	ctx := context.Background()

	testCases := []struct {
		Name   string
		Create func(*testing.T) []transaction.Transaction
	}{
		{
			Name: "create customer",
			Create: func(t *testing.T) []transaction.Transaction {
				return newTransactions()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			err := repository.Create(ctx, func() ([]transaction.Transaction, error) {
				transactions := tc.Create(t)
				return transactions, nil
			})
			assert.Equal(t, err, nil)
		})
	}
}

func TestRepository(t *testing.T) {
	t.Parallel()

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		t.Error("error create mysql connection")
	}

	testutils.SeedCustomer(t, db, "customer-tx-1", "John doe", "66-100-1000")
	testutils.SeedSources(t, db)
	r := NewMysqlTransactionRepository(db)

	t.Run("TestCreate", func(t *testing.T) {
		t.Parallel()
		testCreate(t, r)
	})
}
