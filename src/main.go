package main

import (
	"context"
	"net/http"
	creditLimitPorts "xyz_multifinance/src/internal/creditlimit/ports"
	creditlimitService "xyz_multifinance/src/internal/creditlimit/service"
	customerPorts "xyz_multifinance/src/internal/customer/ports"
	customerService "xyz_multifinance/src/internal/customer/service"
	"xyz_multifinance/src/internal/shared/logs"
	"xyz_multifinance/src/internal/shared/server"
	transactionPort "xyz_multifinance/src/internal/transaction/ports"
	transactionService "xyz_multifinance/src/internal/transaction/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	customerService := customerService.NewApplication(ctx)
	creditLimitService := creditlimitService.NewApplication(ctx)
	transactionService := transactionService.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		customerPorts.HandlerWithOptions(
			customerPorts.NewHttpServer(customerService),
			customerPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []customerPorts.MiddlewareFunc{},
			},
		)

		creditLimitPorts.HandlerWithOptions(
			creditLimitPorts.NewHttpServer(creditLimitService),
			creditLimitPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []creditLimitPorts.MiddlewareFunc{},
			},
		)

		return transactionPort.HandlerWithOptions(
			transactionPort.NewHttpServer(transactionService), transactionPort.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []transactionPort.MiddlewareFunc{},
			},
		)

	})
}
