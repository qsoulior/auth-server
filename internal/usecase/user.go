package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	lowerChars   = `abcdefghijklmnopqrstuvwxyz`
	upperChars   = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	digitChars   = `0123456789`
	specialChars = ` !"#$%&'()*+,-./:;<=>?@[\]^_{|}~`
)

func validateName(name string) error {
	if length := len(name); length < 4 || length > 20 {
		return ErrNameInvalid
	}

	for _, r := range name {
		if !strings.ContainsRune(lowerChars+upperChars+digitChars+"_", r) {
			return ErrNameInvalid
		}
	}

	return nil
}

func validatePassword(password []byte) error {
	if length := len(password); length < 8 || length > 72 {
		return ErrPasswordInvalid
	}

	var lower, upper, digit, special bool

	for _, r := range string(password) {
		switch {
		case strings.ContainsRune(lowerChars, r):
			lower = true
		case strings.ContainsRune(upperChars, r):
			upper = true
		case strings.ContainsRune(digitChars, r):
			digit = true
		case strings.ContainsRune(specialChars, r):
			special = true
		}

		if lower && upper && digit && special {
			return nil
		}
	}

	return ErrPasswordInvalid
}

type UserRepos struct {
	User  repo.User
	Token repo.Token
}

type UserParams struct {
	HashCost int
}

type user struct {
	repos  UserRepos
	params UserParams
}

func NewUser(repos UserRepos, params UserParams) *user {
	return &user{repos, params}
}

func (u *user) hashPassword(password []byte) ([]byte, error) {
	if err := validatePassword(password); err != nil {
		return nil, NewError(err, true)
	}

	hash, err := bcrypt.GenerateFromPassword(password, u.params.HashCost)
	if err != nil {
		return nil, NewError(err, false)
	}

	return hash, nil
}

func (u *user) Create(data entity.User) (*entity.User, error) {
	_, err := u.repos.User.GetByName(context.Background(), data.Name)
	if err == nil {
		return nil, NewError(ErrUserExists, true)
	} else if !errors.Is(err, repo.ErrNoRows) {
		return nil, NewError(err, false)
	}

	if err := validateName(data.Name); err != nil {
		return nil, NewError(err, true)
	}

	hash, err := u.hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	data.Password = hash

	user, err := u.repos.User.Create(context.Background(), data)
	if err != nil {
		return nil, NewError(err, false)
	}

	return user, nil
}

func (u *user) Get(id uuid.UUID) (*entity.User, error) {
	user, err := u.repos.User.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, NewError(ErrUserNotExist, true)
		}
		return nil, NewError(err, false)
	}

	return user, nil
}

func (u *user) Delete(id uuid.UUID) error {
	user, err := u.Get(id)
	if err != nil {
		return err
	}

	if err := u.repos.User.DeleteByID(context.Background(), user.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}

func (u *user) UpdatePassword(id uuid.UUID, password []byte) error {
	user, err := u.Get(id)
	if err != nil {
		return err
	}

	hash, err := u.hashPassword(password)
	if err != nil {
		return err
	}

	if err = u.repos.User.UpdatePassword(context.Background(), user.ID, hash); err != nil {
		return NewError(err, false)
	}

	if err = u.repos.Token.DeleteByUser(context.Background(), user.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}
