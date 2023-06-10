package proxy

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type token struct {
	usecase usecase.Token
}

func NewToken(usecase usecase.Token) *token {
	return &token{usecase}
}

func (t *token) verify(token *entity.RefreshToken, fpData []byte) error {
	fp := usecase.NewFingerprint(fpData, token.UserID)

	if err := fp.Verify(token.Fingerprint); err != nil {
		return usecase.NewError(err, true)
	}

	return nil
}

func (t *token) Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	return t.usecase.Authorize(data, fingerprint)
}

func (t *token) Refresh(id uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.usecase.Get(id)
	if err != nil {
		return "", nil, err
	}

	if err := t.verify(token, fingerprint); err != nil {
		return "", nil, err
	}

	return t.usecase.Refresh(id)
}

func (t *token) Delete(id uuid.UUID, fingerprint []byte) error {
	token, err := t.usecase.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, fingerprint); err != nil {
		return err
	}

	return t.usecase.Delete(id)
}

func (t *token) DeleteAll(id uuid.UUID, fingerprint []byte) error {
	token, err := t.usecase.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, fingerprint); err != nil {
		return err
	}

	return t.usecase.DeleteAll(id)
}
