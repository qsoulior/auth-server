package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/pkg/fingerprint"
	"github.com/qsoulior/auth-server/internal/pkg/hash"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// auth implements Auth interface.
type auth struct {
	jwt jwt.Parser
}

// NewAuth creates a new auth use case.
// It returns pointer to an auth instance.
func NewAuth(jwt jwt.Parser) *auth {
	return &auth{jwt}
}

// Verify verifies user's fingerprint, parses access token,
// and retrieves user ID and roles from it.
// It returns user ID and roles if token is correct and not expired.
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
