package v1

import (
	"encoding/json"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type user struct {
	userUC usecase.User
}

func (u *user) Create(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.userUC.Create(user)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}

func (u *user) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("userID").(uuid.UUID)

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&body)
	if err != nil {
		api.ErrorJSON(w, "decoding error "+err.Error(), http.StatusBadRequest)
		return
	}

	err = u.userUC.UpdatePassword(userID, []byte(body.CurrentPassword), []byte(body.NewPassword))
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}
