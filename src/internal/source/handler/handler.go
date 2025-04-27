package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"xyz_multifinance/src/internal/shared/server/httperr"
	"xyz_multifinance/src/internal/source/usecase"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type HttpHandler struct {
	usecase   usecase.Usecase
	jwtSecret []byte
}

func NewHttpHandler(usecase usecase.Usecase, jwtSecret []byte) HttpHandler {
	return HttpHandler{usecase: usecase, jwtSecret: jwtSecret}
}

func (h HttpHandler) RegisterNewSource(w http.ResponseWriter, r *http.Request) {
	var req RegisterNewSourceJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	if err := h.usecase.Register(r.Context(), usecase.Source{
		ID:       uuid.NewString(),
		Secret:   req.Secret,
		Category: req.Category,
		Name:     req.Name,
		Email:    req.Email,
	}); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, Message{
		Message: "Source created successfully.",
	})
}

func (h HttpHandler) TokenGeneration(w http.ResponseWriter, r *http.Request) {
	var req TokenGenerationJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	source, err := h.usecase.FindByID(r.Context(), req.SourceId)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	if err := source.Validate(req.SourceSecret); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      source.ID,
		"name":     source.Name,
		"category": source.Category,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
}
