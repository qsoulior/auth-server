package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/pkg/fingerprint"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// TokenRepos represents repositories the token use case interacts with.
type TokenRepos struct {
	Token repo.Token
	Role  repo.Role
}

// TokenParams represents parameters for token use case.
type TokenParams struct {
	AccessAge  int
	RefreshAge int
	RefreshCap int
}

// Validate compares parameters with min and max values.
// It returns error if at least one of parameters is invalid.
func (p TokenParams) Validate() error {
	if p.AccessAge < 1 || p.AccessAge > 60 {
		return ErrAccessAgeInvalid
	}
	if p.RefreshAge < 1 || p.RefreshAge > 366 {
		return ErrRefreshAgeInvalid
	}
	if p.RefreshCap < 1 {
		return ErrRefreshCapInvalid
	}
	return nil
}

// token implements Token interface.
type token struct {
	repos  TokenRepos
	params TokenParams
	jwt    jwt.Builder
}

// NewToken validates parameters and creates a new token use case.
// It returns pointer to a token instance or nil if parameters are invalid.
func NewToken(repos TokenRepos, params TokenParams, jwt jwt.Builder) (*token, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &token{repos, params, jwt}, nil
}

// verify compares user's fingerprint with token-related fingerprint.
// It returns nil if fingerprints are equal.
func (t *token) verify(token *entity.RefreshToken, fp []byte) error {
	fpObj := fingerprint.New(token.UserID, fp)
	if err := fpObj.Verify(token.Fingerprint); err != nil {
		return NewError(err, true)
	}

	return nil
}

// create creates new access and refresh tokens using user's fingerprint.
// It returns entity.AccessToken instance
// and pointer to an entity.RefreshToken instance.
func (t *token) create(userID uuid.UUID, fp []byte, session bool) (entity.AccessToken, *entity.RefreshToken, error) {
	// fingerprint
	fpObj := fingerprint.New(userID, fp)
	fpHash, err := fpObj.Hash()
	if err != nil {
		return "", nil, NewError(err, true)
	}

	// refresh token
	rtData := entity.RefreshToken{
		ExpiresAt:   time.Now().AddDate(0, 0, t.params.RefreshAge),
		Fingerprint: fpHash,
		Session:     session,
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

	roleTitles := make([]string, len(roles))
	for i, role := range roles {
		roleTitles[i] = role.Title
	}

	at, err := t.jwt.Build(userID.String(), time.Duration(t.params.AccessAge)*time.Minute, fpHash.HexString(), roleTitles)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	return entity.AccessToken(at), rt, nil
}

// Create creates new access and refresh tokens using user's fingerprint
// and deletes old tokens if total number of tokens is greater than RefreshCap.
// It returns entity.AccessToken instance
// and pointer to an entity.RefreshToken instance.
func (t *token) Create(userID uuid.UUID, fp []byte, session bool) (entity.AccessToken, *entity.RefreshToken, error) {
	tokens, err := t.repos.Token.GetByUser(context.Background(), userID)
	if err != nil {
		return "", nil, NewError(err, false)
	}

	if len(tokens) >= t.params.RefreshCap {
		if err := t.repos.Token.DeleteByID(context.Background(), tokens[0].ID); err != nil {
			return "", nil, NewError(err, false)
		}
	}

	accessToken, refreshToken, err := t.create(userID, fp, session)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

// Refresh verifies user's fingerprint and current refresh token by ID,
// deletes an old refresh token and creates new access and refresh tokens.
// It returns entity.AccessToken instance
// and pointer to an entity.RefreshToken instance.
func (t *token) Refresh(id uuid.UUID, fp []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	token, err := t.Get(id)
	if err != nil {
		return "", nil, err
	}

	if err := t.verify(token, fp); err != nil {
		return "", nil, err
	}

	if err := t.repos.Token.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, NewError(err, false)
	}

	accessToken, refreshToken, err := t.create(token.UserID, fp, token.Session)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

// Get gets a refresh token by ID.
// It returns pointer to an entity.RefreshToken instance
// if id is correct and token isn't expired.
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

// Delete verifies user's fingerprint and current refresh token by ID
// and deletes a refresh token by ID.
// It returns error if id is incorrect or token is expired.
func (t *token) Delete(id uuid.UUID, fp []byte) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, fp); err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByID(context.Background(), token.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}

// DeleteAll verifies user's fingerprint and current refresh token by ID
// and deletes all user refresh tokens.
// It returns error if id is incorrect or token is expired.
func (t *token) DeleteAll(id uuid.UUID, fp []byte) error {
	token, err := t.Get(id)
	if err != nil {
		return err
	}

	if err := t.verify(token, fp); err != nil {
		return err
	}

	if err = t.repos.Token.DeleteByUser(context.Background(), token.UserID); err != nil {
		return NewError(err, false)
	}

	return nil
}
