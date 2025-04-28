package ports

import (
	"net/http"
	"xyz_multifinance/src/internal/creditlimit/app"
	"xyz_multifinance/src/internal/creditlimit/app/command"
	"xyz_multifinance/src/internal/creditlimit/app/query"
	"xyz_multifinance/src/internal/shared/server/httperr"

	"github.com/go-chi/render"
)

type HttpServer struct {
	service app.Application
}

func NewHttpServer(service app.Application) HttpServer {
	return HttpServer{
		service: service,
	}
}
func (h HttpServer) SetInitialCreditLimit(w http.ResponseWriter, r *http.Request, customerId string) {
	request := &SetInitialCreditLimitJSONRequestBody{}
	if err := render.Decode(r, request); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	cmd := command.SetInitialTenorLimit{
		CustomerID: customerId,
		Tenors: []struct {
			MonthRange  int
			LimitAmount float64
		}{},
	}
	for _, tenor := range *request {
		cmd.Tenors = append(cmd.Tenors, struct {
			MonthRange  int
			LimitAmount float64
		}{
			MonthRange:  tenor.MonthRange,
			LimitAmount: tenor.LimitAmount,
		})
	}

	if err := h.service.Commands.SetInitialTenorLimit.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, Message{
		Message: "Initial credit limits created successfully.",
	})

}

func (h HttpServer) GetInitialCreditLimit(w http.ResponseWriter, r *http.Request, customerId string) {
	creditLimit, err := h.service.Queries.GetCreditLimitByCustomerID.Handle(
		r.Context(),
		query.GetCreditLimitByCustomerID{
			CustomerID: customerId},
	)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var response TenorLimits
	for _, t := range creditLimit.Tenors {
		response = append(response, struct {
			LimitAmount float64  "json:\"limitAmount\""
			MonthRange  int      "json:\"monthRange\""
			UsedAmount  *float64 "json:\"usedAmount,omitempty\""
		}{
			LimitAmount: t.LimitAmount,
			MonthRange:  t.MonthRange,
			UsedAmount:  &t.UsedAmount,
		})
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, response)
}
