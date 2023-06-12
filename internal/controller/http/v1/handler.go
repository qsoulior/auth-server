package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

func Handler(user usecase.User, token usecase.Token, logger log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/user/", http.StripPrefix("/user", &UserHandler{user, token, logger}))
	mux.Handle("/token/", http.StripPrefix("/token", &TokenHandler{token, logger}))

	return mux
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

func readAccessToken(r *http.Request) entity.AccessToken {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) < 2 || authorization[0] != "Bearer" {
		return ""
	}

	return authorization[1]
}

func readRefreshToken(r *http.Request) (uuid.UUID, error) {
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

func readFingerprint(r *http.Request) []byte {
	return []byte(fmt.Sprintf("%s : %s : %s : %s", r.Header.Get("Sec-CH-UA"), r.Header.Get("User-Agent"), r.Header.Get("Accept-Language"), r.Header.Get("Upgrade-Insecure-Requests")))
}

func writeSuccess(w http.ResponseWriter) {
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"message": "success",
	})
}

func writeAccessToken(w http.ResponseWriter, token entity.AccessToken) {
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"access_token": token,
	})
}

func writeRefreshToken(w http.ResponseWriter, token *entity.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token.ID.String(),
		Expires:  token.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func deleteRefreshToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "refresh_token",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}