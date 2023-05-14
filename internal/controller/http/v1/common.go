package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

const api = "/api/v1"

func auth(r *http.Request) (entity.AccessToken, error) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) < 2 || authorization[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return entity.AccessToken(authorization[1]), nil
}

func readUser(r *http.Request) (entity.User, error) {
	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		return user, errors.New("decoding error")
	}
	return user, nil
}

func readToken(r *http.Request) (uuid.UUID, error) {
	var (
		data  string
		token uuid.UUID
	)

	if cookie, err := r.Cookie("refresh_token"); err != http.ErrNoCookie {
		data = cookie.Value
	}

	if data == "" {
		return token, errors.New("token is empty")
	}

	return uuid.FromString(data)
}

func writeSuccess(w http.ResponseWriter) {
	w.WriteHeader(200)

	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}

func writeToken(w http.ResponseWriter, access entity.AccessToken, refresh *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh.ID.String(),
		Expires:  refresh.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(200)

	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token": access,
	})
}

func deleteToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "refresh_token",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}
