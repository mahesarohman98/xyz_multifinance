package app

import "xyz_multifinance/src/internal/transaction/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SubmitLoad command.SubmitLoanHandler
}

type Queries struct {
}
