package transaction

import "errors"

type Source struct {
	ID         string
	ExternalID string
}

type Transaction struct {
	ID             string
	ContractNumber string
	CustomerID     string
	Source         Source
	Tenor          int
	OTR            float64
	AmountInterest float64
	AssetName      string
}

// AdminFee lets assume admin fee is calculated as 1% of OTR.
func (t Transaction) AdminFee() float64 {
	return t.OTR * 0.01
}

// TotalBorowed calculates the total borrowed amount
func (t Transaction) TotalBorowed() float64 {
	return t.OTR + t.AdminFee() + t.AmountInterest
}

func (t Transaction) InstallmentAmount() float64 {
	return t.TotalBorowed() / float64(t.Tenor)
}

func NewTransaction(
	iD string,
	contractNumber string,
	customerID string,
	source Source,
	tenor int,
	oTR float64,
	amountInterest float64,
	assetName string,
) (*Transaction, error) {
	if tenor <= 0 {
		return nil, errors.New("tenor must be greater than zero")
	}
	return &Transaction{
		ID:             iD,
		ContractNumber: contractNumber,
		CustomerID:     customerID,
		Source:         source,
		Tenor:          tenor,
		OTR:            oTR,
		AmountInterest: amountInterest,
		AssetName:      assetName,
	}, nil
}
