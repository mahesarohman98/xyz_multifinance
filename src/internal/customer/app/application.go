package app

import (
	"xyz_multifinance/src/internal/customer/app/command"
	"xyz_multifinance/src/internal/customer/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterNewCustomer command.RegisterNewCustomerHandler
}

type Queries struct {
	GetCustomerByID query.GetCustomerByIDHandler
}
