package ports

import (
	"context"
	"net/http"
	"time"
	"xyz_multifinance/src/internal/customer/app"
	"xyz_multifinance/src/internal/customer/app/command"
	"xyz_multifinance/src/internal/customer/app/query"
	"xyz_multifinance/src/internal/shared/server/httperr"
	"xyz_multifinance/src/pkg/dateparser"

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

func (h HttpServer) getCustomerByID(ctx context.Context, customerID string) (*Customer, error) {
	customer, err := h.service.Queries.GetCustomerByID.Handle(ctx, query.GetCustomerByID{
		CustomerByID: customerID,
	})
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:           &customer.ID,
		Nik:          customer.NIK,
		FullName:     customer.Fullname,
		LegalName:    customer.LegalName,
		PlaceOfBirth: customer.PlaceAndDateOfBirth.Place,
		DateOfBirth:  dateparser.UnmarshallToString(customer.PlaceAndDateOfBirth.Date),
		Wage:         customer.Wage,
	}, nil
}

func (h HttpServer) RegisterNewCustomer(w http.ResponseWriter, r *http.Request) {
	request := &Customer{}
	if err := render.Decode(r, request); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	dob, err := dateparser.ParseDate(request.DateOfBirth)
	if err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	cmd := command.RegisterNewCustomer{
		ID:           uuid.NewString(),
		NIK:          request.Nik,
		Fullname:     request.FullName,
		LegalName:    request.LegalName,
		PlaceOfBirth: request.PlaceOfBirth,
		DateOfBirth:  dob,
		Wage:         request.Wage,
		Today:        time.Now(),
	}

	if err := h.service.Commands.RegisterNewCustomer.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	customer, err := h.getCustomerByID(r.Context(), cmd.ID)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, customer)

}

func (h HttpServer) GetCustomerByID(w http.ResponseWriter, r *http.Request, customerId string) {
	customer, err := h.getCustomerByID(r.Context(), customerId)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, customer)
}
