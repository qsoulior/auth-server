package v1

import (
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
)

type token struct {
	userUC  usecase.User
	tokenUC usecase.Token
}

func (t *token) Create(w http.ResponseWriter, r *http.Request) {
	data, err := readUser(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	userID, err := t.userUC.Verify(data)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	accessToken, refreshToken, err := t.tokenUC.Create(userID, fingerprint)
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

	accessToken, refreshToken, err := t.tokenUC.Refresh(token, fingerprint)
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

	err = t.tokenUC.Delete(token, fingerprint)
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

	err = t.tokenUC.DeleteAll(token, fingerprint)
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
