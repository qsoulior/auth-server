package v1

import (
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
)

type token struct {
	usecase usecase.Token
}

func HandleToken(usecase usecase.Token, mux *http.ServeMux) {
	token := &token{usecase}
	mux.Handle(api+"/token/", http.StripPrefix(api+"/token", token))
}

func (t *token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/authorize":
		if r.Method == http.MethodPost {
			t.Authorize(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/refresh":
		if r.Method == http.MethodPost {
			t.Refresh(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/revoke":
		if r.Method == http.MethodPost {
			t.Revoke(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/revoke-all":
		if r.Method == http.MethodPost {
			t.RevokeAll(w, r)
		} else {
			controller.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	default:
		controller.NotFound(w, r)
	}
}

func (t *token) Authorize(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := t.usecase.Authorize(user)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeToken(w, accessToken, refreshToken)
}

func (t *token) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := t.usecase.Refresh(token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		if err == usecase.ErrTokenExpired {
			deleteToken(w)
		}
		return
	}

	writeToken(w, accessToken, refreshToken)
}

func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		if err == usecase.ErrTokenExpired {
			deleteToken(w)
		}
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.usecase.Revoke(token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	deleteToken(w)
	writeSuccess(w)
}

func (t *token) RevokeAll(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.usecase.RevokeAll(token)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	deleteToken(w)
	writeSuccess(w)
}
