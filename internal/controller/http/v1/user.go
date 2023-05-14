package v1

import (
	"encoding/json"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase/proxy"
)

type user struct {
	usecase proxy.User
}

func HandleUser(usecase proxy.User, mux *http.ServeMux) {
	user := &user{usecase}
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
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccess(w)
}

func (u *user) ChangePassword(w http.ResponseWriter, r *http.Request) {
	token, err := auth(r)
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
		controller.ErrorJSON(w, "decoding error"+err.Error(), http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdatePassword(body.NewPassword, body.CurrentPassword, token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccess(w)
}
