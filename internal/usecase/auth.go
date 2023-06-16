package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/pkg/fingerprint"
	"github.com/qsoulior/auth-server/internal/pkg/hash"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type auth struct {
	jwt jwt.Parser
}

func NewAuth(jwt jwt.Parser) *auth {
	return &auth{jwt}
}

func (a *auth) Verify(token entity.AccessToken, fp []byte) (uuid.UUID, []string, error) {
	claims, err := a.jwt.Parse(string(token))
	if err != nil {
		return uuid.UUID{}, nil, NewError(ErrTokenInvalid, true)
	}

	userID, err := uuid.FromString(claims.Subject)
	if err != nil {
		return uuid.UUID{}, nil, NewError(ErrUserIDInvalid, true)
	}

	fpObj := fingerprint.New(userID, fp)
	if err := fpObj.Verify(hash.FromHexString(claims.Fingerprint)); err != nil {
		return uuid.UUID{}, nil, NewError(err, true)
	}

	return userID, claims.Roles, nil
}
