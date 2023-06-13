package v1

import (
	"encoding/json"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type UserHandler struct {
	userUsecase usecase.User
	logger      log.Logger
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		switch r.Method {
		case http.MethodPost:
			u.Create(w, r)
		case http.MethodGet:
			break
		case http.MethodDelete:
			break
		default:
			api.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodGet, http.MethodDelete})
		}
	case "/password":
		if r.Method == http.MethodPut {
			u.UpdatePassword(w, r)
		} else {
			api.MethodNotAllowed(w, r, []string{http.MethodPut})
		}
	default:
		api.NotFound(w, r)
	}
}

func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.userUsecase.Create(user)
	if err != nil {
		api.HandleError(w, err, u.logger, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}

func (u *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	token := readAccessToken(r)
	fingerprint := readFingerprint(r)
	userID, err := u.userUsecase.Authorize(token, fingerprint)
	if err != nil {
		api.HandleError(w, err, u.logger, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusForbidden)
		})
		return
	}

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	d := json.NewDecoder(r.Body)
	err = d.Decode(&body)
	if err != nil {
		api.ErrorJSON(w, "decoding error "+err.Error(), http.StatusBadRequest)
		return
	}

	err = u.userUsecase.UpdatePassword(userID, []byte(body.CurrentPassword), []byte(body.NewPassword))
	if err != nil {
		api.HandleError(w, err, u.logger, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}
