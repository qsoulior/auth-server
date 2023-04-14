package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/rand"
)

type Token struct {
	repo repo.Token
}

func (t *Token) Refresh(userId int) (*entity.Token, error) {
	data, err := rand.GetString(64)
	if err != nil {
		return nil, err
	}

	token := entity.Token{
		Data:      data,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err = t.repo.Create(context.Background(), token, userId)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) RefreshSilent(current entity.Token, userId int) (*entity.Token, error) {
	stored, err := t.repo.GetByUser(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	if current.Data != stored.Data {
		return nil, fmt.Errorf("token is incorrect")
	}
	if current.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}
	newToken, err := t.Refresh(userId)
	return newToken, err
}

func (t *Token) Revoke(userId int) error {
	return t.repo.DeleteByUser(context.Background(), userId)
}

func NewToken(repo repo.Token) *Token {
	return &Token{repo}
}
