package usecase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	token *Token
	repo  repo.User
}

func (u *User) SignUp(user entity.User) error {
	_, err := u.repo.GetByName(context.Background(), user.Name)
	if err == nil {
		return fmt.Errorf("user already exists")
	} else if err != pgx.ErrNoRows {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = u.repo.Create(context.Background(), user)

	return err
}

func (u *User) SignIn(user entity.User) (*entity.Token, error) {
	existing, err := u.repo.GetByName(context.Background(), user.Name)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(user.Password))
	if err != nil {
		return nil, err
	}

	token, err := u.token.Refresh(existing.Id)
	return token, err
}

func (u *User) SignOut(user entity.User) error {
	return u.token.Revoke(user.Id)
}

func NewUser(tu *Token, repo repo.User) *User {
	return &User{tu, repo}
}
