package creditlimit

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestAddTenorLimit(t *testing.T) {
	t.Parallel()

	type args struct {
		monthRange  int
		limitAmount float64
	}

	tests := []struct {
		name            string
		creditLimit     *CreditLimit
		args            args
		wantErr         error
		wantCreditLimit *CreditLimit
	}{
		{
			name: "insert while empty",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors:     []Tenor{},
			},
			args: args{
				monthRange:  1,
				limitAmount: 100,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert duplicate",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  1,
				limitAmount: 100,
			},
			wantErr: ErrDuplicateMonth,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 6",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  6,
				limitAmount: 600,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 2",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  2,
				limitAmount: 200,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 3 error",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  3,
				limitAmount: 200,
			},
			wantErr: ErrLimitMustBeGreater,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 3 error 2",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  3,
				limitAmount: 800,
			},
			wantErr: ErrLimitMustBeLessThanNext,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 3",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  3,
				limitAmount: 300,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "insert 3 (2)",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:  3,
				limitAmount: 300,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			err := tt.creditLimit.AddTenor(tt.args.monthRange, tt.args.limitAmount)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantCreditLimit, tt.creditLimit)
		})
	}
}

func TestDecreaseLimit(t *testing.T) {
	t.Parallel()

	type args struct {
		monthRange   int
		totalBorowed float64
	}

	tests := []struct {
		name            string
		creditLimit     *CreditLimit
		args            args
		wantErr         error
		wantCreditLimit *CreditLimit
	}{
		{
			name: "decrease monthRange 3",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  0,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:   3,
				totalBorowed: 100,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  100,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
		{
			name: "decrease monthRange 3 full",
			creditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  100,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
			args: args{
				monthRange:   3,
				totalBorowed: 200,
			},
			wantErr: nil,
			wantCreditLimit: &CreditLimit{
				CustomerID: "customerid-1",
				Tenors: []Tenor{
					{
						MonthRange:  1,
						LimitAmount: 100,
						UsedAmount:  0,
					},
					{
						MonthRange:  2,
						LimitAmount: 200,
						UsedAmount:  0,
					},
					{
						MonthRange:  3,
						LimitAmount: 300,
						UsedAmount:  300,
					},
					{
						MonthRange:  6,
						LimitAmount: 600,
						UsedAmount:  0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			err := tt.creditLimit.DecreaseLimit(tt.args.monthRange, tt.args.totalBorowed)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantCreditLimit, tt.creditLimit)
		})
	}
}
