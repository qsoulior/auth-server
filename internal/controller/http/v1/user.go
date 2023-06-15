package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type user struct {
	userUsecase usecase.User
	logger      log.Logger
}

func UserMux(userUsecase usecase.User, logger log.Logger) http.Handler {
	u := user{userUsecase, logger}
	auth := AuthMiddleware(userUsecase, logger)
	mux := chi.NewMux()
	mux.Post("/", u.Create)
	mux.With(auth).Put("/password", u.UpdatePassword)
	return mux
}

func (u *user) Create(w http.ResponseWriter, r *http.Request) {
	user, err := readUser(r)
	if err != nil {
		api.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.userUsecase.Create(user)
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

	err = u.userUsecase.UpdatePassword(userID, []byte(body.CurrentPassword), []byte(body.NewPassword))
	if err != nil {
		api.HandleError(err, func(e *usecase.Error) {
			api.ErrorJSON(w, e.Err.Error(), http.StatusBadRequest)
		})
		return
	}

	w.WriteHeader(200)
	writeSuccess(w)
}
