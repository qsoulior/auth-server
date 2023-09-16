package v1

import (
	"encoding/json"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
)

// token represents controllers grouped by token route.
type token struct {
	userUC  usecase.User
	tokenUC usecase.Token
}

// Create reads user data and fingerprint from request, calls User.Verify
// use case to authenticate user and Token.Create use case to create
// new access and refresh tokens.
func (t *token) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Session  bool   `json:"session"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		api.DecodingError(w)
		return
	}

	fingerprint := readFingerprint(r)

	userID, err := t.userUC.Verify(entity.User{Name: data.Name, Password: []byte(data.Password)})
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	accessToken, refreshToken, err := t.tokenUC.Create(userID, fingerprint, data.Session)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	writeRefreshToken(w, refreshToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeAccessToken(w, accessToken)
}

// Refresh reads refresh token and fingerprint from request
// and calls Token.Refresh use case to create new access and refresh tokens.
func (t *token) Refresh(w http.ResponseWriter, r *http.Request) {
	tokenID, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	accessToken, refreshToken, err := t.tokenUC.Refresh(tokenID, fingerprint)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeAccessToken(w, accessToken)
}

// Revoke reads refresh token and fingerprint from request
// and calls Token.Delete to delete refresh token.
func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
	tokenID, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUC.Delete(tokenID, fingerprint)
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
	w.WriteHeader(http.StatusNoContent)
}

// RevokeAll reads refresh token and fingerprint from request
// and calls Token.DeleteAll to delete all user refresh tokens.
func (t *token) RevokeAll(w http.ResponseWriter, r *http.Request) {
	tokenID, err := readRefreshToken(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	fingerprint := readFingerprint(r)

	err = t.tokenUC.DeleteAll(tokenID, fingerprint)
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
	w.WriteHeader(http.StatusNoContent)
}
