package ports

import (
	"net/http"
	"xyz_multifinance/src/internal/shared/server/httperr"
	"xyz_multifinance/src/internal/transaction/app"
	"xyz_multifinance/src/internal/transaction/app/command"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HttpServer struct {
	service app.Application
}

func NewHttpServer(service app.Application) HttpServer {
	return HttpServer{
		service: service,
	}
}

func (h HttpServer) SubmitLoan(w http.ResponseWriter, r *http.Request) {
	request := &SubmitLoanJSONRequestBody{}
	if err := render.Decode(r, request); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	cmd := command.SubmitLoan{
		CustomerID: request.CustomerId,
		Source: struct {
			ID         string
			ExternalID string
		}{
			// TODO: authenticate source
			ID:         "sourceid-1",
			ExternalID: request.ExternalId,
		},
		Tenor: request.Tenor,
		Loans: []command.Loan{},
	}
	for _, loan := range request.Loans {
		cmd.Loans = append(cmd.Loans, command.Loan{
			ID:             uuid.NewString(),
			ContractNumber: loan.ContractNumber,
			OTR:            loan.Otr,
			AmountInterest: loan.AmountInterest,
			AssetName:      loan.AssetName,
		})
	}

	if err := h.service.Commands.SubmitLoad.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, Message{
		Message: "Initial credit limits created successfully.",
	})
}
