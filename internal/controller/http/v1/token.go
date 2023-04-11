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

func (t *token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// create >> signin -> access+refresh
	// update -> access+refresh
	// delete >> signout
}
