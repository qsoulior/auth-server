package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
)

type user struct {
	usecase usecase.User
}

func NewUserController(usecase usecase.User) *user {
	return &user{usecase}
}

func (u *user) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/signup":
		if r.Method == http.MethodPost {
			u.SignUp(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/signin":
		if r.Method == http.MethodPost {
			u.SignIn(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
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

	_, err = u.usecase.SignUp(user)
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
	user, err := readUser(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := u.usecase.SignIn(user)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToken(w, accessToken, refreshToken)
}

func (u *user) ChangePassword(w http.ResponseWriter, r *http.Request) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) < 2 || authorization[0] != "Bearer" {
		controller.ErrorJSON(w, "invalid authorization header", http.StatusBadRequest)
		return
	}
	accessToken := entity.AccessToken(authorization[1])

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&body)
	if err != nil {
		controller.ErrorJSON(w, "decoding error", http.StatusBadRequest)
		return
	}

	err = u.usecase.ChangePassword(body.CurrentPassword, body.NewPassword, accessToken)
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

func readUser(r *http.Request) (*entity.User, error) {
	user := new(entity.User)
	d := json.NewDecoder(r.Body)
	err := d.Decode(user)
	if err != nil {
		return nil, errors.New("decoding error")
	}
	return user, nil
}
