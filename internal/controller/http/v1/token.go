package v1

import (
	"net/http"

	handler "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type TokenHandler struct {
	tokenUsecase usecase.Token
	logger       log.Logger
}

func (t *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/authorize":
		if r.Method == http.MethodPost {
			t.Authorize(w, r)
		} else {
			handler.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/refresh":
		if r.Method == http.MethodPost {
			t.Refresh(w, r)
		} else {
			handler.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/revoke":
		if r.Method == http.MethodPost {
			t.Revoke(w, r)
		} else {
			handler.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	case "/revoke-all":
		if r.Method == http.MethodPost {
			t.RevokeAll(w, r)
		} else {
			handler.MethodNotAllowed(w, r, []string{http.MethodPost})
		}
	default:
		handler.NotFound(w, r)
	}
}

func (t *TokenHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		handler.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	accessToken, refreshToken, err := t.tokenUsecase.Authorize(user, fingerprint)
	if err != nil {
		handler.HandleError(w, err, t.logger, func(e *usecase.Error) {
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeRefreshToken(w, refreshToken)
	w.WriteHeader(200)
	writeAccessToken(w, accessToken)
}

func (t *TokenHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		handler.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	accessToken, refreshToken, err := t.tokenUsecase.Refresh(token, fingerprint)
	if err != nil {
		handler.HandleError(w, err, t.logger, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeRefreshToken(w, refreshToken)
	w.WriteHeader(200)
	writeAccessToken(w, accessToken)
}

func (t *TokenHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		handler.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUsecase.Delete(token, fingerprint)
	if err != nil {
		handler.HandleError(w, err, t.logger, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	deleteRefreshToken(w)
	w.WriteHeader(200)
	writeSuccess(w)
}

func (t *TokenHandler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		handler.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUsecase.DeleteAll(token, fingerprint)
	if err != nil {
		handler.HandleError(w, err, t.logger, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			handler.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	deleteRefreshToken(w)
	w.WriteHeader(200)
	writeSuccess(w)
}
