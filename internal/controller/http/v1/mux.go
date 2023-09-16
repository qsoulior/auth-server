// Package v1 provides structures and functions to implement HTTP controllers.
package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// Mux creates a new mux and mounts controllers.
// It returns pointer to a chi.Mux instance.
func Mux(userUC usecase.User, tokenUC usecase.Token, authUC usecase.Auth, logger log.Logger) http.Handler {
	user := user{userUC}
	token := token{userUC, tokenUC}
	auth := AuthMiddleware(authUC, logger)

	mux := chi.NewMux()
	mux.Route("/", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/", user.Create)
			r.With(auth).Get("/", user.Get)
			r.With(auth).Delete("/", user.Delete)
			r.With(auth).Put("/password", user.UpdatePassword)
		})
		r.Route("/token", func(r chi.Router) {
			r.Post("/", token.Create)
			r.Post("/refresh", token.Refresh)
			r.Post("/revoke", token.Revoke)
			r.Post("/revoke-all", token.RevokeAll)
		})
	})

	return mux
}

// readAccessToken reads access token from request's Authorization header.
// It returns access token string or empty string if header is invalid.
func readAccessToken(r *http.Request) (entity.AccessToken, error) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) < 2 || authorization[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return entity.AccessToken(authorization[1]), nil
}

// readRefreshToken reads refresh token from request's cookie.
// It returns error if cookie is empty.
func readRefreshToken(r *http.Request) (uuid.UUID, error) {
	var (
		data  string
		token uuid.UUID
	)

	if cookie, err := r.Cookie("refresh_token"); err != http.ErrNoCookie {
		data = cookie.Value
	}

	if data == "" {
		return token, errors.New("token is empty")
	}

	return uuid.FromString(data)
}

// readFingerprint reads headers from request and creates fingerprint using them.
// It returns fingerprint byte slice.
func readFingerprint(r *http.Request) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s", r.Header.Get("Sec-CH-UA"), r.Header.Get("User-Agent"), r.Header.Get("Accept-Language"), r.Header.Get("Upgrade-Insecure-Requests")))
}

// writeAccessToken writes an access token to response body.
func writeAccessToken(w http.ResponseWriter, token entity.AccessToken) {
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token": token,
	})
}

// writeRefreshToken writes a refresh token to response cookie.
func writeRefreshToken(w http.ResponseWriter, token *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Path:     "/v1/token",
		Value:    token.ID.String(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	if !token.Session {
		cookie.Expires = token.ExpiresAt
	}

	http.SetCookie(w, cookie)
}

// deleteRefreshToken writes an expired refresh token to response cookie.
func deleteRefreshToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
}
