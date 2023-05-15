package v1

import (
	"errors"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type token struct {
	usecase usecase.Token
	logger  log.Logger
}

func HandleToken(usecase usecase.Token, mux *http.ServeMux, logger log.Logger) {
	token := &token{usecase, logger}
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
		var e *usecase.Error
		if errors.As(err, &e) && e.External {
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
			return
		}
		controller.InternalServerError(w)
		t.logger.Error("%s", err)
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
		var e *usecase.Error
		if errors.As(err, &e) && e.External {
			if errors.Is(e.Err, usecase.ErrTokenExpired) {
				deleteToken(w)
			}
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
			return
		}
		controller.InternalServerError(w)
		t.logger.Error("%s", err)
		return
	}

	writeToken(w, accessToken, refreshToken)
}

func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := readToken(r)
	if err != nil {
		controller.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.usecase.Revoke(token)
	if err != nil {
		var e *usecase.Error
		if errors.As(err, &e) && e.External {
			if errors.Is(e.Err, usecase.ErrTokenExpired) {
				deleteToken(w)
			}
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
			return
		}
		controller.InternalServerError(w)
		t.logger.Error("%s", err)
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
		var e *usecase.Error
		if errors.As(err, &e) && e.External {
			if errors.Is(e.Err, usecase.ErrTokenExpired) {
				deleteToken(w)
			}
			controller.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
			return
		}
		controller.InternalServerError(w)
		t.logger.Error("%s", err)
		return
	}

	deleteToken(w)
	writeSuccess(w)
}
