package v1

import (
	"context"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

func AuthMiddleware(userUsecase usecase.User, logger log.Logger) api.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := readAccessToken(r)
			fingerprint := readFingerprint(r)
			userID, err := userUsecase.Authorize(token, fingerprint)
			if err != nil {
				api.HandleError(err, func(e *usecase.Error) {
					api.ErrorJSON(w, e.Err.Error(), http.StatusForbidden)
				})
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
