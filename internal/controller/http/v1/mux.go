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

func readUser(r *http.Request) (*entity.User, error) {
	user := new(entity.User)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func readAccessToken(r *http.Request) (entity.AccessToken, error) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) < 2 || authorization[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return entity.AccessToken(authorization[1]), nil
}

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

func readFingerprint(r *http.Request) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s", r.Header.Get("Sec-CH-UA"), r.Header.Get("User-Agent"), r.Header.Get("Accept-Language"), r.Header.Get("Upgrade-Insecure-Requests")))
}

func writeAccessToken(w http.ResponseWriter, token entity.AccessToken) {
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token": token,
	})
}

func writeRefreshToken(w http.ResponseWriter, token *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Path:     "/v1/token",
		Value:    token.ID.String(),
		Expires:  token.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
}

func deleteRefreshToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "refresh_token",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}
