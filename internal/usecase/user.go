package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrNameInvalid       = errors.New("name is invalid")
	ErrPasswordInvalid   = errors.New("password is invalid")
	ErrPasswordIncorrect = errors.New("password is incorrect")
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

func validatePassword(password string) error {
	if length := len(password); length < 8 || length > 72 {
		return ErrPasswordInvalid
	}

	var lower, upper, digit, special bool

	for _, r := range password {
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

type User struct {
	token    *Token
	repo     repo.User
	hashCost int
}

func NewUser(tu *Token, repo repo.User, hashCost int) *User {
	return &User{tu, repo, hashCost}
}

func (u *User) SignUp(user entity.User) error {
	_, err := u.repo.GetByName(context.Background(), user.Name)
	if err == nil {
		return ErrUserExists
	} else if err != repo.ErrUserNotExist {
		return err
	}

	if err := validateName(user.Name); err != nil {
		return err
	}

	if err := validatePassword(user.Password); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), u.hashCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	err = u.repo.Create(context.Background(), user)

	return err
}

func (u *User) SignIn(user entity.User) (*entity.AccessToken, *entity.RefreshToken, error) {
	existing, err := u.repo.GetByName(context.Background(), user.Name)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(user.Password))
	if err != nil {
		return nil, nil, ErrPasswordIncorrect
	}

	return u.token.Refresh(existing.ID)
}
