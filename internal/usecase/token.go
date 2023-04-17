package usecase

import (
	"context"
	"errors"
	"time"

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
}

func (t *Token) Refresh(userID int) (*entity.Token, error) {
	data, err := rand.GetString(64)
	if err != nil {
		return nil, err
	}

	token := entity.Token{
		Data:      data,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err = t.repo.Create(context.Background(), token, userID)
	if err != nil {
		return nil, err
	}

	return &token, nil
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

func (t *Token) RefreshSilent(current string, userID int) (*entity.Token, error) {
	if err := t.checkToken(current, userID); err != nil {
		return nil, err
	}

	newToken, err := t.Refresh(userID)
	// TODO: generate access token
	return newToken, err
}

func (t *Token) Revoke(current string, userID int) error {
	if err := t.checkToken(current, userID); err != nil {
		return err
	}

	return t.repo.DeleteByUser(context.Background(), userID)
}

func NewToken(repo repo.Token) *Token {
	return &Token{repo}
}
