package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

var (
	ErrTokenIncorrect = errors.New("token is incorrect")
	ErrTokenExpired   = errors.New("token is expired")
)

type TokenParams struct {
	Issuer    string
	Algorithm string
	Key       any
}

type token struct {
	repo    repo.Token
	builder jwt.Builder
}

func NewToken(repo repo.Token, params TokenParams) (*token, error) {
	builder, err := jwt.NewBuilder(params.Issuer, params.Algorithm, params.Key)
	if err != nil {
		return nil, err
	}

	return &token{repo, builder}, nil
}

func (t *token) Refresh(userID int) (entity.AccessToken, *entity.RefreshToken, error) {
	data, err := uuid.New()
	if err != nil {
		return "", nil, err
	}

	refreshToken := entity.RefreshToken{
		Data:      data,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err = t.repo.Create(context.Background(), refreshToken, userID)
	if err != nil {
		return "", nil, err
	}

	token, err := t.builder.Build(strconv.Itoa(userID))
	if err != nil {
		return "", nil, err
	}

	accessToken := entity.AccessToken(token)

	return accessToken, &refreshToken, nil
}

func (t *token) getToken(data uuid.UUID) (*entity.RefreshToken, error) {
	token, err := t.repo.GetByData(context.Background(), data)
	if err != nil {
		return nil, ErrTokenIncorrect
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	return token, nil
}

func (t *token) RefreshSilent(data uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.getToken(data)
	if err != nil {
		return "", nil, err
	}

	if err := t.repo.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, err
	}

	return t.Refresh(token.UserID)
}

func (t *token) Revoke(data uuid.UUID) error {
	token, err := t.getToken(data)
	if err != nil {
		return err
	}

	return t.repo.DeleteByID(context.Background(), token.ID)
}
