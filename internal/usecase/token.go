package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/pkg/fingerprint"
	"github.com/qsoulior/auth-server/internal/pkg/hash"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type TokenRepos struct {
	Token repo.Token
	Role  repo.Role
}

type TokenJWT struct {
	jwt.Builder
	jwt.Parser
}

type TokenParams struct {
	AccessAge  int
	RefreshAge int
	RefreshCap int
}

type token struct {
	repos  TokenRepos
	jwt    TokenJWT
	params TokenParams
}

func NewToken(repos TokenRepos, jwt TokenJWT, params TokenParams) *token {
	return &token{repos, jwt, params}
}

func (t *token) verify(token *entity.RefreshToken, userData []byte) error {
	fp := fingerprint.New(token.UserID, userData)
	if err := fp.Verify(token.Fingerprint); err != nil {
		return NewError(err, true)
	}

	return nil
}

func (t *token) create(userID uuid.UUID, userData []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	// fingerprint
	fp := fingerprint.New(userID, userData)
	fpHash, err := fp.Hash()
	if err != nil {
		return "", nil, NewError(err, true)
	}

	// refresh token
	rtData := entity.RefreshToken{
		ExpiresAt:   time.Now().AddDate(0, 0, t.params.RefreshAge),
		Fingerprint: fpHash,
		UserID:      userID,
	}

	rt, err := t.repos.Token.Create(context.Background(), rtData)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	// access token
	roles, err := t.repos.Role.GetByUser(context.Background(), userID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	rolesID := make([]string, len(roles))
	for i, role := range roles {
		rolesID[i] = role.ID.String()
	}

	at, err := t.jwt.Build(userID.String(), time.Duration(t.params.AccessAge)*time.Minute, fpHash.HexString(), rolesID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	return entity.AccessToken(at), rt, nil
}

func (t *token) Create(userID uuid.UUID, userData []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	tokens, err := t.repos.Token.GetByUser(context.Background(), userID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	if len(tokens) >= t.params.RefreshCap {
		if err := t.repos.Token.DeleteByID(context.Background(), tokens[0].ID); err != nil {
			return "", nil, NewError(err, false)
		}
	}

	accessToken, refreshToken, err := t.create(userID, userData)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *token) Refresh(id uuid.UUID, userData []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.Get(id)
	if err != nil {
		return "", nil, err
	}

	if err := t.verify(token, userData); err != nil {
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

func (t *token) Authorize(token entity.AccessToken, userData []byte) (uuid.UUID, error) {
	var userID uuid.UUID
	claims, err := t.jwt.Parse(string(token))
	if err != nil {
		return userID, NewError(err, true)
	}

	userID, err = uuid.FromString(claims.Subject)
	if err != nil {
		return userID, NewError(ErrUserIDInvalid, true)
	}

	fp := fingerprint.New(userID, userData)
	if err := fp.Verify(hash.FromHexString(claims.Fingerprint)); err != nil {
		return userID, NewError(err, true)
	}

	return userID, nil
}

func (t *token) Delete(id uuid.UUID, userData []byte) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, userData); err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByID(context.Background(), token.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}

func (t *token) DeleteAll(id uuid.UUID, userData []byte) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, userData); err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByUser(context.Background(), token.UserID); err != nil {
		return NewError(err, false)
	}

	return nil
}
