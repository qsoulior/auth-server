package v1

import (
	"net/http"

	"github.com/qsoulior/auth-server/internal/usecase"
)

type user struct {
	usecase *usecase.User
}

func NewUser(u *usecase.User) *user {
	return &user{u}
}

func (u *user) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// create
	// get
	// update
	// delete
}
