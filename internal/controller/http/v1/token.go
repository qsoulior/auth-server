package v1

import (
	"net/http"

	"github.com/qsoulior/auth-server/internal/usecase"
)

type token struct {
	usecase *usecase.Token
}

func NewToken(u *usecase.Token) *token {
	return &token{u}
}

func (t *token) Refresh(w http.ResponseWriter, r *http.Request) {
}

func (t *token) Revoke(w http.ResponseWriter, r *http.Request) {
}
