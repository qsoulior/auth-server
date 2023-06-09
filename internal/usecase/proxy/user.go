package proxy

import (
	"encoding/hex"

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

func (u *user) verifyToken(token entity.AccessToken, fingerprint []byte) (uuid.UUID, error) {
	var userID uuid.UUID
	sub, fp, err := u.parser.Parse(token)
	if err != nil {
		return userID, err
	}

	userID, err = uuid.FromString(sub)
	if err != nil {
		return userID, usecase.ErrUserIDInvalid
	}

	fingerprint, err = usecase.HashFingerprint(userID, fingerprint)
	if err != nil {
		return userID, err
	}

	if hex.EncodeToString(fingerprint) != fp {
		return userID, usecase.ErrFingerprintIncorrect
	}

	return userID, nil
}

func (u *user) verifyPassword(id uuid.UUID, password []byte) error {
	user, err := u.usecase.Get(id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, password); err != nil {
		return usecase.ErrPasswordIncorrect
	}

	return nil
}

func (u *user) Create(data entity.User) (*entity.User, error) {
	return u.usecase.Create(data)
}

func (u *user) Get(token entity.AccessToken, fingerprint []byte) (*entity.User, error) {
	id, err := u.verifyToken(token, fingerprint)
	if err != nil {
		return nil, usecase.NewError(err, true)
	}

	return u.usecase.Get(id)
}

func (u *user) Delete(password []byte, token entity.AccessToken, fingerprint []byte) error {
	id, err := u.verifyToken(token, fingerprint)
	if err != nil {
		return usecase.NewError(err, true)
	}

	if err = u.verifyPassword(id, password); err != nil {
		return usecase.NewError(err, true)
	}

	return u.usecase.Delete(id)
}

func (u *user) UpdatePassword(newPassword []byte, password []byte, token entity.AccessToken, fingerprint []byte) error {
	id, err := u.verifyToken(token, fingerprint)
	if err != nil {
		return usecase.NewError(err, true)
	}

	if err = u.verifyPassword(id, password); err != nil {
		return usecase.NewError(err, true)
	}

	return u.usecase.UpdatePassword(id, newPassword)
}
