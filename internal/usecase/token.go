package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/rand"
)

var (
	ErrTokenIncorrect = errors.New("token is incorrect")
	ErrTokenExpired   = errors.New("token is expired")
)

type Token struct {
	repo repo.Token
	key  any
}

func NewToken(repo repo.Token, key any) *Token {
	return &Token{repo, key}
}

func (t *Token) Refresh(userID int) (*entity.AccessToken, *entity.RefreshToken, error) {
	data, err := rand.GetString(64)
	if err != nil {
		return nil, nil, err
	}

	refreshToken := entity.RefreshToken{
		Data:      data,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err = t.repo.Create(context.Background(), refreshToken, userID)
	if err != nil {
		return nil, nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Minute),
	})

	tokenStr, err := token.SignedString(t.key)
	if err != nil {
		return nil, nil, err
	}

	accessToken := entity.AccessToken(tokenStr)

	return &accessToken, &refreshToken, nil
}

func (t *Token) checkToken(current string, userID int) error {
	stored, err := t.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return err
	}

	if current != stored.Data {
		return ErrTokenIncorrect
	}

	if stored.ExpiresAt.Before(time.Now()) {
		return ErrTokenExpired
	}
	return nil
}

func (t *Token) RefreshSilent(current string, userID int) (*entity.AccessToken, *entity.RefreshToken, error) {
	if err := t.checkToken(current, userID); err != nil {
		return nil, nil, err
	}

	return t.Refresh(userID)
}

func (t *Token) Revoke(current string, userID int) error {
	if err := t.checkToken(current, userID); err != nil {
		return err
	}

	return t.repo.DeleteByUser(context.Background(), userID)
}
