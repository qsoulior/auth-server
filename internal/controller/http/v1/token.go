package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

type token struct {
	userUsecase  usecase.User
	tokenUsecase usecase.Token
	logger       log.Logger
}

func TokenMux(userUsecase usecase.User, tokenUsecase usecase.Token, logger log.Logger) http.Handler {
	t := token{userUsecase, tokenUsecase, logger}
	mux := chi.NewMux()
	mux.Post("/", t.Create)
	mux.Post("/refresh", t.Refresh)
	mux.Post("/revoke", t.Revoke)
	mux.Post("/revoke-all", t.RevokeAll)
	return mux
}

func (t *token) Create(w http.ResponseWriter, r *http.Request) {
	data, err := readUser(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	userID, err := t.userUsecase.Authenticate(data)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	accessToken, refreshToken, err := t.tokenUsecase.Create(userID, fingerprint)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeRefreshToken(w, refreshToken)
	w.WriteHeader(200)
	writeAccessToken(w, accessToken)
}

func (t *token) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	accessToken, refreshToken, err := t.tokenUsecase.Refresh(token, fingerprint)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeRefreshToken(w, refreshToken)
	w.WriteHeader(200)
	writeAccessToken(w, accessToken)
}

func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUsecase.Delete(token, fingerprint)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	deleteRefreshToken(w)
	w.WriteHeader(200)
	writeSuccess(w)
}

func (t *token) RevokeAll(w http.ResponseWriter, r *http.Request) {
	token, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUsecase.DeleteAll(token, fingerprint)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			if e.Err == usecase.ErrTokenExpired {
				deleteRefreshToken(w)
			}
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	deleteRefreshToken(w)
	w.WriteHeader(200)
	writeSuccess(w)
}
