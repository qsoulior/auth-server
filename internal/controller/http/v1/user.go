package v1

import (
	"encoding/json"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type user struct {
	usecase *usecase.User
	logger  log.Logger
}

func NewUserController(usecase *usecase.User, logger log.Logger) *user {
	return &user{usecase, logger}
}

func (u *user) SignUp(w http.ResponseWriter, r *http.Request) {
	controllerName := "sign up"
	address := r.RemoteAddr

	if r.Method != http.MethodPost {
		err := controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		err := controller.UnsupportedMediaType(w, r, controller.ContentType)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		err := controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	err = u.usecase.SignUp(user)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	w.WriteHeader(200)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}

func (u *user) SignIn(w http.ResponseWriter, r *http.Request) {
	controllerName := "sign in"
	address := r.RemoteAddr

	if r.Method != http.MethodPost {
		err := controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		err := controller.UnsupportedMediaType(w, r, controller.ContentType)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		err := controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	accessToken, refreshToken, err := u.usecase.SignIn(user)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		u.logger.Error(controller.Error{err, controllerName, address})
		return
	}

	writeToken(w, accessToken, refreshToken)
}
