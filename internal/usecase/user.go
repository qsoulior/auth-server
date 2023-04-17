package usecase

import (
	"context"
	"errors"
	"regexp"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrNameRegexp        = errors.New("name must be at least 4 characters...")
	ErrPasswordRegexp    = errors.New("password must be at least 8 characters...")
	ErrPasswordIncorrect = errors.New("password is incorrect")
)

var (
	NameRegexp     = regexp.MustCompile("^[A-Za-z0-9]{4,30}$")
	PasswordRegexp = regexp.MustCompile("^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[#?!@$%^&*_(),.+-]).{8,30}$")
)

type User struct {
	token *Token
	repo  repo.User
}

func (u *User) SignUp(user entity.User) error {
	_, err := u.repo.GetByName(context.Background(), user.Name)
	if err == nil {
		return ErrUserExists
	} else if err != repo.ErrUserNotExist {
		return err
	}

	if !NameRegexp.MatchString(user.Name) {
		return ErrNameRegexp
	}

	if !PasswordRegexp.MatchString(user.Password) {
		return ErrPasswordRegexp
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
		return nil, ErrPasswordIncorrect
	}

	return u.token.Refresh(existing.ID)
}

func NewUser(tu *Token, repo repo.User) *User {
	return &User{tu, repo}
}
