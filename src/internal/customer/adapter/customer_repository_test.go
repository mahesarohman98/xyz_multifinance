package adapter

import (
	"context"
	"testing"
	"time"
	"xyz_multifinance/src/internal/customer/domain/customer"
	"xyz_multifinance/src/internal/shared/mysql"

	"gopkg.in/go-playground/assert.v1"
)

var factory = customer.NewFactory(customer.FactoryConfig{
	WageLimit:  3500000,
	MinimumAge: 20,
})

type args struct {
	id           string
	nik          string
	fullname     string
	legalName    string
	placeOfBirth string
	dateOfBirth  time.Time
	wage         float64
	photoURL     string
	kTPURL       string
	processDate  time.Time
}

func newValidCustomer() args {
	return args{
		id:           "customerid-1",
		nik:          "666-999-111",
		fullname:     "john doe",
		legalName:    "john doe",
		placeOfBirth: "tangerang",
		dateOfBirth:  time.Date(1998, time.April, 1, 1, 0, 0, 0, time.UTC),
		wage:         3500000,
		photoURL:     "https://imgurl.com/photo",
		kTPURL:       "https://imgurl.com/ktp",
		processDate:  time.Date(2025, time.April, 1, 1, 0, 0, 0, time.UTC),
	}
}

func newCustomer(t *testing.T) *customer.Customer {
	args := newValidCustomer()
	c, err := factory.RegisterNewCustomer(
		args.id,
		args.nik,
		args.fullname,
		args.legalName,
		args.placeOfBirth,
		args.dateOfBirth,
		args.wage,
		args.photoURL,
		args.kTPURL,
		args.processDate,
	)
	assert.Equal(t, err, nil)
	return c
}

func TestRepository(t *testing.T) {
	t.Parallel()

	db, err := mysql.NewMySQLConnection()
	if err != nil {
		t.Error("error create mysql connection")
	}
	r := NewMysqlCustomerRepository(db)

	t.Run("TestCreate", func(t *testing.T) {
		t.Parallel()
		testCreate(t, r)
	})

}

func testCreate(t *testing.T, repository customer.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name   string
		Create func(*testing.T) *customer.Customer
	}{
		{
			Name: "create customer",
			Create: func(t *testing.T) *customer.Customer {
				return newCustomer(t)
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
