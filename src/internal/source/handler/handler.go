package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"xyz_multifinance/src/internal/shared/server/httperr"
	"xyz_multifinance/src/internal/source/model"
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

func (h HttpHandler) tokenGeneration(ctx context.Context, source *model.Source, passwordSecret string) (TokenResponse, error) {
	if err := source.Validate(passwordSecret); err != nil {
		return TokenResponse{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      source.ID,
		"name":     source.Name,
		"category": source.Category,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}

func (h HttpHandler) RegisterNewSource(w http.ResponseWriter, r *http.Request) {
	var req RegisterNewSourceJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	source, err := h.usecase.Register(r.Context(), usecase.Source{
		ID:       uuid.NewString(),
		Secret:   req.Secret,
		Category: req.Category,
		Name:     req.Name,
		Email:    req.Email,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	token, err := h.tokenGeneration(r.Context(), source, req.Secret)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Respond(w, r, token)
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
	token, err := h.tokenGeneration(r.Context(), source, req.SourceSecret)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Respond(w, r, token)
}
