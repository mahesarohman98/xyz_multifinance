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

	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	customerService := customerService.NewApplication(ctx)
	creditLimitService := creditlimitService.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		customerPorts.HandlerWithOptions(
			customerPorts.NewHttpServer(customerService),
			customerPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []customerPorts.MiddlewareFunc{},
			},
		)

		return creditLimitPorts.HandlerWithOptions(
			creditLimitPorts.NewHttpServer(creditLimitService),
			creditLimitPorts.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []creditLimitPorts.MiddlewareFunc{},
			})
	})
}
