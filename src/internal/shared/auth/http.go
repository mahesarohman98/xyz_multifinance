package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"xyz_multifinance/src/internal/shared/server/httperr"

	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

var JWTSecret = []byte("your-super-secret-key")

type contextKey string

const (
	SourceIDKey contextKey = "sourceId"
	CategoryKey contextKey = "category"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, ok := ctx.Value(BearerAuthScopes).([]string)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httperr.Unauthorised("missing auth header", fmt.Errorf("missing auth header"), w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httperr.Unauthorised("invalid auth format", fmt.Errorf("invalid auth format"), w, r)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			httperr.Unauthorised("invalid token", fmt.Errorf("invalid token"), w, r)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		sourceID := claims["sub"].(string)
		scopes := claims["category"].(string)

		ctx = context.WithValue(r.Context(), SourceIDKey, sourceID)
		ctx = context.WithValue(ctx, CategoryKey, scopes)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireCategory(requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scope, ok := r.Context().Value(CategoryKey).(string)
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if scope == requiredScope {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func SourceIDFromCtx(ctx context.Context) (string, error) {
	sourceID, ok := ctx.Value(SourceIDKey).(string)
	if !ok {
		return "", errors.New("sourceID not found")
	}
	return sourceID, nil
}
