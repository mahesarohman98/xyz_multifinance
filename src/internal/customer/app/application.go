package app

import "xyz_multifinance/src/internal/customer/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterNewCustomer command.RegisterNewCustomerHandler
}

type Queries struct {
}
