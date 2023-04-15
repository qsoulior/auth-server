package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotExist      = errors.New("user does not exist")
	ErrPasswordShort     = errors.New("password must be at least 8 characters")
	ErrPasswordIncorrect = errors.New("password is incorrect")
)

type User struct {
	token *Token
	repo  repo.User
}

func (u *User) SignUp(user entity.User) error {
	_, err := u.repo.GetByName(context.Background(), user.Name)
	if err == nil {
		return ErrUserExists
	} else if err != pgx.ErrNoRows {
		return err
	}

	if len(user.Password) < 8 {
		return ErrPasswordShort
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
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotExist
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(user.Password))
	if err != nil {
		return nil, ErrPasswordIncorrect
	}

	token, err := u.token.Refresh(existing.Id)
	// TODO: generate access token
	return token, err
}

func (u *User) SignOut(user entity.User) error {
	return u.token.Revoke(user.Id)
}

func NewUser(tu *Token, repo repo.User) *User {
	return &User{tu, repo}
}
