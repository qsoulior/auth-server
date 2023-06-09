package usecase

import (
	"crypto/sha256"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(id uuid.UUID) (*entity.User, error)
	Delete(id uuid.UUID) error
	UpdatePassword(id uuid.UUID, password []byte) error
}

type Token interface {
	Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Refresh(id uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Revoke(id uuid.UUID, fingerprint []byte) error
	RevokeAll(id uuid.UUID, fingerprint []byte) error
}

func HashFingerprint(userID uuid.UUID, fingerprint []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(append(fingerprint, userID[:]...))
	if err != nil {
		return nil, ErrFingerprintInvalid
	}

	return h.Sum(nil), nil
}
