package v1

import (
	"encoding/json"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
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

	var body struct {
		UserID int    `json:"user_id"`
		Token  string `json:"token"`
	}

	d := json.NewDecoder(r.Body)
	err := d.Decode(&body)
	if err != nil {
		err := controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		t.logger.Error(controller.NewError(err, controllerName, address))
		return
	}

	if r.Method == http.MethodPost {
		token, err := t.usecase.RefreshSilent(body.Token, body.UserID)
		if err != nil {
			err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
			t.logger.Error(controller.NewError(err, controllerName, address))
			return
		}

		w.WriteHeader(200)
		e := json.NewEncoder(w)
		e.Encode(token)
		return
	}

	err = t.usecase.Revoke(body.Token, body.UserID)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		t.logger.Error(controller.NewError(err, controllerName, address))
		return
	}

	w.WriteHeader(200)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}
