package ports

import (
	"net/http"
	"time"
	"xyz_multifinance/src/internal/customer/app"
	"xyz_multifinance/src/internal/customer/app/command"
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
	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, Message{
		Message: "Customer created successfully.",
	})

}
