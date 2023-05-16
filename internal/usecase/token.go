package usecase

import (
	"context"
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
}

type token struct {
	userRepo  repo.User
	tokenRepo repo.Token
	builder   jwt.Builder

	accessAge  int
	refreshAge int
}

func NewToken(userRepo repo.User, tokenRepo repo.Token, builder jwt.Builder, params TokenParams) *token {
	return &token{userRepo, tokenRepo, builder, params.AccessAge, params.RefreshAge}
}

func (t *token) getByID(id uuid.UUID) (*entity.RefreshToken, error) {
	token, err := t.tokenRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, ErrTokenIncorrect
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	return token, nil
}

func (t *token) create(userID uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error) {
	data := entity.RefreshToken{
		ExpiresAt: time.Now().AddDate(0, 0, t.refreshAge),
		UserID:    userID,
	}

	refreshToken, err := t.tokenRepo.Create(context.Background(), data)
	if err != nil {
		return "", nil, err
	}

	accessToken, err := t.builder.Build(userID.String(), time.Duration(t.accessAge)*time.Minute)
	if err != nil {
		return "", nil, err
	}

	return entity.AccessToken(accessToken), refreshToken, nil
}

func (t *token) Authorize(data entity.User) (entity.AccessToken, *entity.RefreshToken, error) {
	const fn = "Authorize"

	user, err := t.userRepo.GetByName(context.Background(), data.Name)
	if err != nil {
		if err == repo.ErrUserNotExist {
			return "", nil, tokenError(fn, err, true)
		}
		return "", nil, tokenError(fn, err, false)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return "", nil, tokenError(fn, ErrPasswordIncorrect, true)
	}

	accessToken, refreshToken, err := t.create(user.ID)
	if err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Refresh(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error) {
	const fn = "Refresh"

	token, err := t.getByID(id)
	if err != nil {
		return "", nil, tokenError(fn, err, true)
	}

	if err := t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	accessToken, refreshToken, err := t.create(token.UserID)
	if err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Revoke(id uuid.UUID) error {
	const fn = "Revoke"

	token, err := t.getByID(id)
	if err != nil {
		return tokenError(fn, err, true)
	}

	if err = t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return tokenError(fn, err, false)
	}

	return nil
}

func (t *token) RevokeAll(id uuid.UUID) error {
	const fn = "RevokeAll"

	token, err := t.getByID(id)
	if err != nil {
		return tokenError(fn, err, true)
	}

	if err = t.tokenRepo.DeleteByUser(context.Background(), token.UserID); err != nil {
		return tokenError(fn, err, false)
	}

	return nil
}
