package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
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

type UserParams struct {
	Issuer    string
	Algorithm string
	Key       any
	HashCost  int
}

type user struct {
	token    Token
	repo     repo.User
	parser   jwt.Parser
	hashCost int
}

func NewUser(tu Token, repo repo.User, params UserParams) (*user, error) {
	parser, err := jwt.NewParser(params.Issuer, params.Algorithm, params.Key)
	if err != nil {
		return nil, err
	}
	return &user{tu, repo, parser, params.HashCost}, nil
}

func (u *user) auth(token entity.AccessToken) (int, error) {
	sub, err := u.parser.Parse(string(token))
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		return 0, errors.New("user id is invalid")
	}
	return id, nil
}

func (u *user) hashPassword(password string) ([]byte, error) {
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	return bcrypt.GenerateFromPassword([]byte(password), u.hashCost)
}

func (u *user) SignUp(user *entity.User) error {
	_, err := u.repo.GetByName(context.Background(), user.Name)
	if err == nil {
		return ErrUserExists
	} else if err != repo.ErrUserNotExist {
		return err
	}

	if err := validateName(user.Name); err != nil {
		return err
	}

	hash, err := u.hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return u.repo.Create(context.Background(), user)
}

func (u *user) SignIn(user *entity.User) (entity.AccessToken, *entity.RefreshToken, error) {
	ex, err := u.repo.GetByName(context.Background(), user.Name)
	if err != nil {
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(ex.Password), []byte(user.Password))
	if err != nil {
		return "", nil, ErrPasswordIncorrect
	}

	return u.token.Refresh(ex.ID)
}

func (u *user) ChangePassword(password string, newPassword string, accessToken entity.AccessToken) error {
	id, err := u.auth(accessToken)
	if err != nil {
		return err
	}

	user, err := u.repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ErrPasswordIncorrect
	}

	hash, err := u.hashPassword(newPassword)
	if err != nil {
		return err
	}

	return u.repo.UpdatePassword(context.Background(), id, string(hash))
}
