package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenRepos struct {
	User  repo.User
	Token repo.Token
	Role  repo.Role
}

type TokenParams struct {
	AccessAge  int
	RefreshAge int
	RefreshCap int
}

type token struct {
	repos   TokenRepos
	params  TokenParams
	builder jwt.Builder
}

func NewToken(repos TokenRepos, params TokenParams, builder jwt.Builder) *token {
	return &token{repos, params, builder}
}

func (t *token) create(userID uuid.UUID, fpData []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	fp := NewFingerprint(fpData, userID)
	fpBytes, err := fp.Hash()
	if err != nil {
		return "", nil, NewError(err, true)
	}
	fpString := hex.EncodeToString(fpBytes)

	data := entity.RefreshToken{
		ExpiresAt:   time.Now().AddDate(0, 0, t.params.RefreshAge),
		Fingerprint: fpBytes,
		UserID:      userID,
	}

	refreshToken, err := t.repos.Token.Create(context.Background(), data)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	roles, err := t.repos.Role.GetByUser(context.Background(), userID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	rolesID := make([]string, len(roles))
	for i, role := range roles {
		rolesID[i] = role.ID.String()
	}

	accessToken, err := t.builder.Build(userID.String(), time.Duration(t.params.AccessAge)*time.Minute, fpString, rolesID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	user, err := t.repos.User.GetByName(context.Background(), data.Name)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return "", nil, NewError(ErrUserNotExist, true)
		}
		return "", nil, NewError(err, false)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, data.Password); err != nil {
		return "", nil, NewError(ErrPasswordIncorrect, true)
	}

	tokens, err := t.repos.Token.GetByUser(context.Background(), user.ID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	if len(tokens) >= t.params.RefreshCap {
		if err := t.repos.Token.DeleteByID(context.Background(), tokens[0].ID); err != nil {
			return "", nil, NewError(err, false)
		}
	}

	accessToken, refreshToken, err := t.create(user.ID, fingerprint)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *token) Refresh(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.Get(id)
	if err != nil {
		return "", nil, err
	}

	if err := t.repos.Token.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, NewError(err, false)
	}

	accessToken, refreshToken, err := t.create(token.UserID, token.Fingerprint)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *token) Get(id uuid.UUID) (*entity.RefreshToken, error) {
	token, err := t.repos.Token.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, NewError(ErrTokenIncorrect, true)
		}
		return nil, NewError(err, false)
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, NewError(ErrTokenExpired, true)
	}
	return token, nil
}

func (t *token) Delete(id uuid.UUID) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByID(context.Background(), token.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}

func (t *token) DeleteAll(id uuid.UUID) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByUser(context.Background(), token.UserID); err != nil {
		return NewError(err, false)
	}

	return nil
}
