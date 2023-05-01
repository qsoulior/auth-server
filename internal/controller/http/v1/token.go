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
	usecase *usecase.Token
	logger  log.Logger
}

func NewTokenController(usecase *usecase.Token, logger log.Logger) *token {
	return &token{usecase, logger}
}

func (t *token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	controllerName := "token"
	address := r.RemoteAddr
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		err := controller.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodDelete})
		t.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		err := controller.UnsupportedMediaType(w, r, controller.ContentType)
		t.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	token, err := readToken(r)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		t.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	if r.Method == http.MethodPost {
		accessToken, refreshToken, err := t.usecase.RefreshSilent(token)
		if err != nil {
			err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
			t.logger.Error(controller.Error{err, controllerName, address})
			return
		}

		writeToken(w, accessToken, refreshToken)
		return
	}

	err = t.usecase.Revoke(token)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		t.logger.Error(controller.Error{err, controllerName, address})
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

func writeToken(w http.ResponseWriter, access *entity.AccessToken, refresh *entity.RefreshToken) {
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
