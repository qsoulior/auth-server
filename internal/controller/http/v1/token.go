package v1

import (
	"encoding/json"
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
		t.logger.Error(controller.NewError(err, controllerName, address))
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		err := controller.UnsupportedMediaType(w, r, controller.ContentType)
		t.logger.Error(controller.NewError(err, controllerName, address))
		return
	}

	var data string
	if cookie, err := r.Cookie("refresh_token"); err != http.ErrNoCookie {
		data = cookie.Value
	} else {
		var body struct {
			Value string `json:"refresh_token"`
		}

		d := json.NewDecoder(r.Body)
		err := d.Decode(&body)
		if err != nil {
			err := controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
			t.logger.Error(controller.NewError(err, controllerName, address))
			return
		}
		data = body.Value
	}

	token, err := uuid.FromString(data)
	if err != nil {
		return
	}

	if r.Method == http.MethodPost {
		accessToken, refreshToken, err := t.usecase.RefreshSilent(token)
		if err != nil {
			err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
			t.logger.Error(controller.NewError(err, controllerName, address))
			return
		}

		writeToken(w, accessToken, refreshToken)
		return
	}

	err = t.usecase.Revoke(token)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		t.logger.Error(controller.NewError(err, controllerName, address))
		return
	}
	deleteToken(w)
}

func writeToken(w http.ResponseWriter, accessToken *entity.AccessToken, refreshToken *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Data.String(),
		Expires:  refreshToken.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)

	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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
