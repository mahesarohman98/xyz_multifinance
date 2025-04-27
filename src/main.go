package main

import (
	"context"
	"net/http"
	customerPort "xyz_multifinance/src/internal/customer/ports"
	customerService "xyz_multifinance/src/internal/customer/service"
	"xyz_multifinance/src/internal/shared/logs"
	"xyz_multifinance/src/internal/shared/server"

	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	customerService := customerService.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return customerPort.HandlerWithOptions(
			customerPort.NewHttpServer(customerService),
			customerPort.ChiServerOptions{
				BaseRouter:  router,
				Middlewares: []customerPort.MiddlewareFunc{},
			},
		)
	})
}
