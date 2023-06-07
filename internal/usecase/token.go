package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/jwt"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenParams struct {
	AccessAge  int
	RefreshAge int
	RefreshCap int
}

type token struct {
	userRepo  repo.User
	tokenRepo repo.Token
	builder   jwt.Builder

	accessAge  int
	refreshAge int
	refreshCap int
}

func NewToken(userRepo repo.User, tokenRepo repo.Token, builder jwt.Builder, params TokenParams) *token {
	return &token{userRepo, tokenRepo, builder, params.AccessAge, params.RefreshAge, params.RefreshCap}
}

func (t *token) getByID(id uuid.UUID) (*entity.RefreshToken, error) {
	token, err := t.tokenRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, ErrTokenIncorrect
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	return token, nil
}

func (t *token) hashFingerprint(userID uuid.UUID, fingerprint []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(append(fingerprint, userID[:]...))
	if err != nil {
		return nil, ErrFingerprintInvalid
	}

	return h.Sum(nil), nil
}

func (t *token) verifyFingerprint(token *entity.RefreshToken, fingerprint []byte) error {
	fp, err := t.hashFingerprint(token.UserID, fingerprint)
	if err != nil {
		return err
	}

	if hex.EncodeToString(fp) != token.Fingerprint {
		return ErrFingerprintIncorrect
	}
	return nil
}

func (t *token) create(userID uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	fp, err := t.hashFingerprint(userID, fingerprint)
	if err != nil {
		return "", nil, err
	}

	data := entity.RefreshToken{
		ExpiresAt:   time.Now().AddDate(0, 0, t.refreshAge),
		Fingerprint: hex.EncodeToString(fp),
		UserID:      userID,
	}

	refreshToken, err := t.tokenRepo.Create(context.Background(), data)
	if err != nil {
		return "", nil, err
	}

	accessToken, err := t.builder.Build(userID.String(), time.Duration(t.accessAge)*time.Minute)
	if err != nil {
		return "", nil, err
	}

	return entity.AccessToken(accessToken), refreshToken, nil
}

func (t *token) Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	const fn = "Authorize"

	user, err := t.userRepo.GetByName(context.Background(), data.Name)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotExist) {
			return "", nil, tokenError(fn, err, true)
		}
		return "", nil, tokenError(fn, err, false)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return "", nil, tokenError(fn, ErrPasswordIncorrect, true)
	}

	tokens, err := t.tokenRepo.GetByUser(context.Background(), user.ID)
	if err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	if len(tokens) >= t.refreshCap {
		if err := t.tokenRepo.DeleteByID(context.Background(), tokens[0].ID); err != nil {
			return "", nil, tokenError(fn, err, false)
		}
	}

	accessToken, refreshToken, err := t.create(user.ID, fingerprint)
	if err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Refresh(id uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error) {
	const fn = "Refresh"

	token, err := t.getByID(id)
	if err != nil {
		return "", nil, tokenError(fn, err, true)
	}

	if err := t.verifyFingerprint(token, fingerprint); err != nil {
		return "", nil, tokenError(fn, err, true)
	}

	if err := t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	accessToken, refreshToken, err := t.create(token.UserID, fingerprint)
	if err != nil {
		return "", nil, tokenError(fn, err, false)
	}

	return accessToken, refreshToken, nil
}

func (t *token) Revoke(id uuid.UUID, fingerprint []byte) error {
	const fn = "Revoke"

	token, err := t.getByID(id)
	if err != nil {
		return tokenError(fn, err, true)
	}

	if err := t.verifyFingerprint(token, fingerprint); err != nil {
		return tokenError(fn, err, true)
	}

	if err = t.tokenRepo.DeleteByID(context.Background(), token.ID); err != nil {
		return tokenError(fn, err, false)
	}

	return nil
}

func (t *token) RevokeAll(id uuid.UUID, fingerprint []byte) error {
	const fn = "RevokeAll"

	token, err := t.getByID(id)
	if err != nil {
		return tokenError(fn, err, true)
	}

	if err := t.verifyFingerprint(token, fingerprint); err != nil {
		return tokenError(fn, err, true)
	}

	if err = t.tokenRepo.DeleteByUser(context.Background(), token.UserID); err != nil {
		return tokenError(fn, err, false)
	}

	return nil
}
