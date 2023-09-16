package v1

import (
	"encoding/json"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// user represents controllers grouped by user route.
type user struct {
	userUC usecase.User
}

// Create reads user from request
// and calls User.Create to create a new user.
func (u *user) Create(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	d := json.NewDecoder(r.Body)
	err := d.Decode(&user)
	if err != nil {
		api.DecodingError(w)
		return
	}

	_, err = u.userUC.Create(user)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get gets user ID from request's context
// and calls User.Get to get user data by ID.
func (u *user) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("userID").(uuid.UUID)

	user, err := u.userUC.Get(userID)
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusNotFound)
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(map[string]any{
		"id":       user.ID,
		"username": user.Name,
	})
}

// Delete gets user ID from request's context and current password
// from request's body, then calls User.Delete to delete user by ID.
func (u *user) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("userID").(uuid.UUID)

	var body struct {
		Password string `json:"password"`
	}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&body)
	if err != nil {
		api.DecodingError(w)
		return
	}

	err = u.userUC.Delete(userID, []byte(body.Password))
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePassword gets user ID from request's context and passwords
// from request's body, then calls User.UpdatePassword to
// update user's password by ID.
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
		api.DecodingError(w)
		return
	}

	err = u.userUC.UpdatePassword(userID, []byte(body.CurrentPassword), []byte(body.NewPassword))
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
