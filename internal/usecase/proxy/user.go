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

func (u *user) verifyToken(token entity.AccessToken, fpData []byte) (uuid.UUID, error) {
	var userID uuid.UUID
	claims, err := u.parser.Parse(token)
	if err != nil {
		return userID, err
	}

	userID, err = uuid.FromString(claims.Subject)
	if err != nil {
		return userID, usecase.ErrUserIDInvalid
	}

	fp := usecase.NewFingerprint(fpData, userID)
	fpBytes, _ := hex.DecodeString(claims.Fingerprint)

	if err := fp.Verify(fpBytes); err != nil {
		return userID, usecase.NewError(err, true)
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

func (u *user) Delete(currentPwd []byte, token entity.AccessToken, fingerprint []byte) error {
	id, err := u.verifyToken(token, fingerprint)
	if err != nil {
		return usecase.NewError(err, true)
	}

	if err = u.verifyPassword(id, currentPwd); err != nil {
		return usecase.NewError(err, true)
	}

	return u.usecase.Delete(id)
}

func (u *user) UpdatePassword(newPwd []byte, currentPwd []byte, token entity.AccessToken, fingerprint []byte) error {
	id, err := u.verifyToken(token, fingerprint)
	if err != nil {
		return usecase.NewError(err, true)
	}

	if err = u.verifyPassword(id, currentPwd); err != nil {
		return usecase.NewError(err, true)
	}

	return u.usecase.UpdatePassword(id, newPwd)
}
