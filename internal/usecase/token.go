package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenParams struct {
	AccessAge  int
	RefreshAge int
	RefreshCap int
}

type token struct {
	userRepo  repo.User
	tokenRepo repo.Token
	builder   jwt.Builder

	accessAge  int
	refreshAge int
	refreshCap int
}

func NewToken(userRepo repo.User, tokenRepo repo.Token, builder jwt.Builder, params TokenParams) *token {
	return &token{userRepo, tokenRepo, builder, params.AccessAge, params.RefreshAge, params.RefreshCap}
}

func (t *token) create(userID uuid.UUID, fpData []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	fp := NewFingerprint(fpData, userID)
	fpBytes, err := fp.Hash()
	if err != nil {
		return "", nil, NewError(err, true)
	}
	fpString := hex.EncodeToString(fpBytes)

	data := entity.RefreshToken{
		ExpiresAt:   time.Now().AddDate(0, 0, t.refreshAge),
		Fingerprint: fpBytes,
		UserID:      userID,
	}

	refreshToken, err := t.tokenRepo.Create(context.Background(), data)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	accessToken, err := t.builder.Build(userID.String(), time.Duration(t.accessAge)*time.Minute, fpString, []string{})
	if err != nil {
		return "", nil, NewError(err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	user, err := t.userRepo.GetByName(context.Background(), data.Name)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return "", nil, NewError(ErrUserNotExist, true)
		}
		return "", nil, NewError(err, false)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, data.Password); err != nil {
		return "", nil, NewError(ErrPasswordIncorrect, true)
	}

	tokens, err := t.tokenRepo.GetByUser(context.Background(), user.ID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	if len(tokens) >= t.refreshCap {
		if err := t.tokenRepo.DeleteByID(context.Background(), tokens[0].ID); err != nil {
			return "", nil, NewError(err, false)
		}
	}

	accessToken, refreshToken, err := t.create(user.ID, fingerprint)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *token) Refresh(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.Get(id)
	if err != nil {
		return "", nil, err
	}

	if err := t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, NewError(err, false)
	}

	accessToken, refreshToken, err := t.create(token.UserID, token.Fingerprint)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *token) Get(id uuid.UUID) (*entity.RefreshToken, error) {
	token, err := t.tokenRepo.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, NewError(ErrTokenIncorrect, true)
		}
		return nil, NewError(err, false)
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, NewError(ErrTokenExpired, true)
	}
	return token, nil
}

func (t *token) Delete(id uuid.UUID) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err = t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}

func (t *token) DeleteAll(id uuid.UUID) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err = t.tokenRepo.DeleteByUser(context.Background(), token.UserID); err != nil {
		return NewError(err, false)
	}

	return nil
}
