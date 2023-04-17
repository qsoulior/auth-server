package v1

import (
	"encoding/json"
	"fmt"
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

func NewUser(usecase *usecase.User, logger log.Logger) *user {
	return &user{usecase, logger}
}

type SignUpError struct {
	Err     error
	Address string
}

func NewSignUpError(err error, address string) *SignUpError {
	return &SignUpError{err, address}
}

func (err *SignUpError) Error() string {
	return fmt.Sprintf("sign up: %s (%s)", err.Err, err.Address)
}

func (u *user) SignUp(w http.ResponseWriter, r *http.Request) {
	address := r.RemoteAddr

	if r.Method != http.MethodPost {
		err := controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		u.logger.Error(NewSignUpError(err, address))
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		err := controller.UnsupportedMediaType(w, r, controller.ContentType)
		u.logger.Error(NewSignUpError(err, address))
		return
	}

	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		err := controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		u.logger.Error(NewSignUpError(err, address))
		return
	}

	err = u.usecase.SignUp(user)
	if err != nil {
		err := controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		u.logger.Error(NewSignUpError(err, address))
		return
	}

	w.WriteHeader(200)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}

func (u *user) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		return
	}

	w.Header().Set("Content-Type", controller.ContentType)
	if r.Header.Get("Content-Type") != controller.ContentType {
		controller.UnsupportedMediaType(w, r, controller.ContentType)
		return
	}

	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		return
	}

	token, err := u.usecase.SignIn(user)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	e := json.NewEncoder(w)
	e.Encode(token)
}
