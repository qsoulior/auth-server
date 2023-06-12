package v1

import (
	"encoding/json"
	"net/http"

	handler "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type UserHandler struct {
	userUsecase  usecase.User
	tokenUsecase usecase.Token
	logger       log.Logger
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		switch r.Method {
		case http.MethodPost:
			u.SignUp(w, r)
		case http.MethodGet:
			break
		case http.MethodDelete:
			break
		default:
			handler.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodGet, http.MethodDelete})
		}
	case "/password":
		if r.Method == http.MethodPut {
			u.ChangePassword(w, r)
		} else {
			handler.MethodNotAllowed(w, r, []string{http.MethodPut})
		}
	default:
		handler.NotFound(w, r)
	}
}

func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		handler.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.userUsecase.Create(user)
	if err != nil {
		handler.HandleError(w, err, u.logger, func(e *usecase.Error) {
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	token := readAccessToken(r)
	fingerprint := readFingerprint(r)
	userID, err := u.tokenUsecase.VerifyAccess(token, fingerprint)
	if err != nil {
		handler.HandleError(w, err, u.logger, func(e *usecase.Error) {
			handler.ErrorJSON(w, e.Err.Error(), http.StatusForbidden)
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
		handler.ErrorJSON(w, "decoding error "+err.Error(), http.StatusBadRequest)
		return
	}

	err = u.userUsecase.UpdatePassword(userID, []byte(body.CurrentPassword), []byte(body.NewPassword))
	if err != nil {
		handler.HandleError(w, err, u.logger, func(e *usecase.Error) {
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}
