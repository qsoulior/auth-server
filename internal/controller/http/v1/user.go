package v1

import (
	"encoding/json"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/internal/usecase/proxy"
	"github.com/qsoulior/auth-server/pkg/log"
)

type user struct {
	usecase proxy.User
	logger  log.Logger
}

func HandleUser(usecase proxy.User, mux *http.ServeMux, logger log.Logger) {
	user := &user{usecase, logger}
	mux.Handle(api+"/user/", http.StripPrefix(api+"/user", user))
}

func (u *user) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			controller.MethodNotAllowed(w, r, []string{http.MethodPost, http.MethodGet, http.MethodDelete})
		}
	case "/password":
		if r.Method == http.MethodPut {
			u.ChangePassword(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPut})
		}
	default:
		controller.NotFound(w, r)
	}
}

func (u *user) SignUp(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.usecase.Create(user)
	if err != nil {
		controller.HandleError(w, err, u.logger, func(e *usecase.Error) {
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeSuccess(w)
}

func (u *user) ChangePassword(w http.ResponseWriter, r *http.Request) {
	token, err := readAuth(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	d := json.NewDecoder(r.Body)
	err = d.Decode(&body)
	if err != nil {
		controller.ErrorJSON(w, "decoding error "+err.Error(), http.StatusBadRequest)
		return
	}

	fingerprint := readFingerprint(r)
	err = u.usecase.UpdatePassword([]byte(body.NewPassword), []byte(body.CurrentPassword), token, fingerprint)
	if err != nil {
		controller.HandleError(w, err, u.logger, func(e *usecase.Error) {
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeSuccess(w)
}
