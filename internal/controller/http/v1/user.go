package v1

import (
	"encoding/json"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
)

type user struct {
	usecase *usecase.User
}

func NewUser(usecase *usecase.User) *user {
	return &user{usecase}
}

func (u *user) SignUp(w http.ResponseWriter, r *http.Request) {
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

	err = u.usecase.SignUp(user)
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
