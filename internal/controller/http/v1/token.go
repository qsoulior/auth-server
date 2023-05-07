package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type token struct {
	usecase usecase.Token
	logger  log.Logger
}

func NewTokenController(usecase usecase.Token, logger log.Logger) *token {
	return &token{usecase, logger}
}

func (t *token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.Refresh(w, r)
	case http.MethodDelete:
		t.Revoke(w, r)
	default:
		controller.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodDelete})
	}
}

func (t *token) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := t.usecase.RefreshSilent(token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToken(w, accessToken, refreshToken)
}

func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.usecase.Revoke(token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	deleteToken(w)
}

func readToken(r *http.Request) (uuid.UUID, error) {
	var (
		data  string
		token uuid.UUID
	)

	if cookie, err := r.Cookie("refresh_token"); err != http.ErrNoCookie {
		data = cookie.Value
	} else {
		var body struct {
			Value string `json:"refresh_token"`
		}

		d := json.NewDecoder(r.Body)
		err := d.Decode(&body)
		if err != nil {
			return token, errors.New("decoding error")
		}
		data = body.Value
	}

	if data == "" {
		return token, errors.New("token is empty")
	}

	return uuid.FromString(data)
}

func writeToken(w http.ResponseWriter, access entity.AccessToken, refresh *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Data.String(),
		Expires:  refresh.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)

	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func deleteToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "refresh_token",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)

	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}
