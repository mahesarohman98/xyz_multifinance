package creditlimit

import (
	"errors"
)

type Tenor struct {
	MonthRange  int
	LimitAmount float64
	UsedAmount  float64
}

type CreditLimit struct {
	CustomerID string
	Tenors     []Tenor
}

type Factory struct {
}

func NewFactory() Factory {
	return Factory{}
}

func (f *Factory) MustNewCreditLimit(customerID string) *CreditLimit {
	return &CreditLimit{
		CustomerID: customerID,
		Tenors:     []Tenor{},
	}
}

var (
	ErrDuplicateMonth          = errors.New("duplicate tenor month range")
	ErrLimitMustBeGreater      = errors.New("limit must be greather than previous limit")
	ErrLimitMustBeLessThanNext = errors.New("limit must be less than next limit")
)

// AddTenor adds a new tenor while maintak to insert tenor.
// return error if month range duplicate or limit less then previous tenor
func (c *CreditLimit) AddTenor(monthRange int, limitAmount float64) error {
	insertIndex := 0
	isLongest := true
	for i, tenor := range c.Tenors {
		insertIndex = i
		if tenor.MonthRange == monthRange {
			return ErrDuplicateMonth
		} else if tenor.MonthRange > monthRange {
			isLongest = false
			break
		}
	}
	if isLongest {
		insertIndex += 1
	}
	if insertIndex >= len(c.Tenors) {
		if insertIndex-1 > 0 && c.Tenors[insertIndex-1].LimitAmount < limitAmount {
			return ErrLimitMustBeGreater
		}
		c.Tenors = append(c.Tenors, Tenor{
			MonthRange:  monthRange,
			LimitAmount: limitAmount,
			UsedAmount:  0,
		})
		return nil
	}
	if insertIndex-1 > 0 && c.Tenors[insertIndex-1].LimitAmount >= limitAmount {
		return ErrLimitMustBeGreater
	}
	if insertIndex < len(c.Tenors) && c.Tenors[insertIndex].LimitAmount <= limitAmount {
		return ErrLimitMustBeLessThanNext
	}
	c.Tenors = append(c.Tenors[:insertIndex+1], c.Tenors[insertIndex:]...)
	c.Tenors[insertIndex] = Tenor{
		MonthRange:  monthRange,
		LimitAmount: limitAmount,
		UsedAmount:  0,
	}

	return nil
}
