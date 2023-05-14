package proxy

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	usecase usecase.User
	parser  jwt.Parser
}

func NewUser(usecase usecase.User, parser jwt.Parser) *user {
	return &user{usecase, parser}
}

func (u *user) verifyToken(token entity.AccessToken) (uuid.UUID, error) {
	var id uuid.UUID
	sub, err := u.parser.Parse(string(token))
	if err != nil {
		return id, err
	}

	id, err = uuid.FromString(sub)
	if err != nil {
		return id, ErrIDInvalid
	}
	return id, nil
}

func (u *user) verifyPassword(id uuid.UUID, password string) error {
	user, err := u.usecase.Get(id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return usecase.ErrPasswordIncorrect
	}

	return nil
}

func (u *user) Create(data entity.User) (*entity.User, error) {
	const fn = "Create"

	return u.usecase.Create(data)
}

func (u *user) Get(token entity.AccessToken) (*entity.User, error) {
	const fn = "Get"

	id, err := u.verifyToken(token)
	if err != nil {
		return nil, userError(fn, err)
	}

	return u.usecase.Get(id)
}

func (u *user) Delete(password string, token entity.AccessToken) error {
	const fn = "Delete"

	id, err := u.verifyToken(token)
	if err != nil {
		return userError(fn, err)
	}

	if err = u.verifyPassword(id, password); err != nil {
		return userError(fn, err)
	}

	return u.usecase.Delete(id)
}

func (u *user) UpdatePassword(newPassword string, password string, token entity.AccessToken) error {
	const fn = "UpdatePassword"

	id, err := u.verifyToken(token)
	if err != nil {
		return userError(fn, err)
	}

	if err = u.verifyPassword(id, password); err != nil {
		return userError(fn, err)
	}

	return u.usecase.UpdatePassword(id, newPassword)
}
