package app

import (
	"xyz_multifinance/src/internal/creditlimit/app/command"
	"xyz_multifinance/src/internal/creditlimit/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SetInitialTenorLimit command.SetInitialTenorLimitHandler
	DecreaseLimit        command.DecreaseLimitHandler
}

type Queries struct {
	GetTotalUsedByCustomerAndTenor query.GetTotalUsedByCustomerAndTenorHandler
	GetCreditLimitByCustomerID     query.GetCreditLimitByCustomerIDHandler
}
