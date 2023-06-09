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

type UserParams struct {
	HashCost int
}

type user struct {
	userRepo  repo.User
	tokenRepo repo.Token

	hashCost int
}

func NewUser(userRepo repo.User, tokenRepo repo.Token, params UserParams) *user {
	return &user{userRepo, tokenRepo, params.HashCost}
}

func (u *user) hashPassword(password []byte) ([]byte, error) {
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(password, u.hashCost)
	if err != nil {
		return nil, ErrPasswordInvalid
	}

	return hash, nil
}

func (u *user) Create(data entity.User) (*entity.User, error) {
	_, err := u.userRepo.GetByName(context.Background(), data.Name)
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
		return nil, NewError(err, true)
	}

	data.Password = hash

	user, err := u.userRepo.Create(context.Background(), data)
	if err != nil {
		return nil, NewError(err, false)
	}

	return user, nil
}

func (u *user) Get(id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, NewError(ErrUserNotExist, true)
		}
		return nil, NewError(err, false)
	}

	return user, nil
}

func (u *user) Delete(id uuid.UUID) error {
	if err := u.userRepo.DeleteByID(context.Background(), id); err != nil {
		return NewError(err, false)
	}

	return nil
}

func (u *user) UpdatePassword(id uuid.UUID, password []byte) error {
	user, err := u.userRepo.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return NewError(ErrUserNotExist, true)
		}
		return NewError(err, false)
	}

	hash, err := u.hashPassword(password)
	if err != nil {
		return NewError(err, true)
	}

	if err = u.userRepo.UpdatePassword(context.Background(), user.ID, hash); err != nil {
		return NewError(err, false)
	}

	if err = u.tokenRepo.DeleteByUser(context.Background(), user.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}
