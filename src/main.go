package main

import (
	"context"
	"net/http"
	creditLimitPorts "xyz_multifinance/src/internal/creditlimit/ports"
	creditlimitService "xyz_multifinance/src/internal/creditlimit/service"
	customerPorts "xyz_multifinance/src/internal/customer/ports"
	customerService "xyz_multifinance/src/internal/customer/service"
	"xyz_multifinance/src/internal/shared/auth"
	"xyz_multifinance/src/internal/shared/logs"
	"xyz_multifinance/src/internal/shared/server"
	"xyz_multifinance/src/internal/source/handler"
	"xyz_multifinance/src/internal/source/service"
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
	sourceService := service.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		customerPorts.HandlerWithOptions(
			customerPorts.NewHttpServer(customerService),
			customerPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []customerPorts.MiddlewareFunc{auth.AuthMiddleware},
			},
		)

		creditLimitPorts.HandlerWithOptions(
			creditLimitPorts.NewHttpServer(creditLimitService),
			creditLimitPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []creditLimitPorts.MiddlewareFunc{auth.AuthMiddleware},
			},
		)

		transactionPort.HandlerWithOptions(
			transactionPort.NewHttpServer(transactionService), transactionPort.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []transactionPort.MiddlewareFunc{auth.AuthMiddleware},
			},
		)

		return handler.HandlerWithOptions(
			handler.NewHttpHandler(sourceService, auth.JWTSecret), handler.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []handler.MiddlewareFunc{},
			},
		)

	})
}
