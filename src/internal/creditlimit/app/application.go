package app

import "xyz_multifinance/src/internal/creditlimit/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SetInitialTenorLimit command.SetInitialTenorLimitHandler
}

type Queries struct {
}
