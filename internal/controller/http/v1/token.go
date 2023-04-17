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

func NewToken(usecase *usecase.Token, logger log.Logger) *token {
	return &token{usecase, logger}
}

func (t *token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		controller.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodDelete})
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		controller.UnsupportedMediaType(w, r, controller.ContentType)
		return
	}

	var body struct {
		UserID int    `json:"user_id"`
		Token  string `json:"token"`
	}

	d := json.NewDecoder(r.Body)
	err := d.Decode(&body)
	if err != nil {
		controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		token, err := t.usecase.RefreshSilent(body.Token, body.UserID)
		if err != nil {
			controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(200)
		e := json.NewEncoder(w)
		e.Encode(token)
		return
	}

	err = t.usecase.Revoke(body.Token, body.UserID)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}
