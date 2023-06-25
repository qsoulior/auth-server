package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	lowerChars   = `abcdefghijklmnopqrstuvwxyz`
	upperChars   = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	digitChars   = `0123456789`
	specialChars = ` !"#$%&'()*+,-./:;<=>?@[\]^_{|}~`
)

// validateName returns error if name is invalid.
func validateName(name string) error {
	if length := len(name); length < 4 || length > 20 {
		return NewError(ErrNameInvalid, true)
	}

	for _, r := range name {
		if !strings.ContainsRune(lowerChars+upperChars+digitChars+"_", r) {
			return NewError(ErrNameInvalid, true)
		}
	}

	return nil
}

// validatePassword returns error if password is invalid.
func validatePassword(password []byte) error {
	if length := len(password); length < 8 || length > 72 {
		return NewError(ErrPasswordInvalid, true)
	}

	var lower, upper, digit, special bool

	for _, r := range string(password) {
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

	return NewError(ErrPasswordInvalid, true)
}

// hashPassword validates password and hashes it using bcrypt algorithm.
// It returns nil if password is invalid or hashing failed.
func hashPassword(password []byte, hashCost int) ([]byte, error) {
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(password, hashCost)
	if err != nil {
		return nil, NewError(err, false)
	}

	return hash, nil
}

// verifyPassword compares hashedPassword with password.
// It returns nil if passwords are equal.
func verifyPassword(hashedPassword []byte, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return NewError(ErrPasswordIncorrect, true)
	}

	return nil
}

// UserRepos represents repositories the user use case interacts with.
type UserRepos struct {
	User repo.User
}

// UserParams represents parameters for user use case.
type UserParams struct {
	HashCost int
}

// Validate compares parameters with min and max values.
// It returns error if at least one of parameters is invalid.
func (p UserParams) Validate() error {
	if p.HashCost < bcrypt.MinCost || p.HashCost > bcrypt.MaxCost {
		return ErrHashCostInvalid
	}
	return nil
}

// user implements User interface.
type user struct {
	repos  UserRepos
	params UserParams
}

// NewUser validates parameters and creates a new user use case.
// It returns pointer to an user instance or nil if parameters are invalid.
func NewUser(repos UserRepos, params UserParams) (*user, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &user{repos, params}, nil
}

// Create validates data and creates a new user.
// It returns pointer to an entity.User instance or nil if an error occurred.
func (u *user) Create(data entity.User) (*entity.User, error) {
	_, err := u.repos.User.GetByName(context.Background(), data.Name)
	if err == nil {
		return nil, NewError(ErrUserExists, true)
	} else if !errors.Is(err, repo.ErrNoRows) {
		return nil, NewError(err, false)
	}

	if err := validateName(data.Name); err != nil {
		return nil, err
	}

	hash, err := hashPassword(data.Password, u.params.HashCost)
	if err != nil {
		return nil, err
	}

	data.Password = hash

	user, err := u.repos.User.Create(context.Background(), data)
	if err != nil {
		return nil, NewError(err, false)
	}

	return user, nil
}

// Get gets a user by ID.
// It returns pointer to an entity.User instance or nil if an error occurred.
func (u *user) Get(id uuid.UUID) (*entity.User, error) {
	user, err := u.repos.User.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, NewError(ErrUserNotExist, true)
		}
		return nil, NewError(err, false)
	}

	return user, nil
}

// Verify verifies user's name and password
// and is used in authentication process.
// It returns user ID if name and password are correct
// or empty UUID if an error occurred.
func (u *user) Verify(data entity.User) (uuid.UUID, error) {
	user, err := u.repos.User.GetByName(context.Background(), data.Name)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return uuid.UUID{}, NewError(ErrUserNotExist, true)
		}
		return uuid.UUID{}, NewError(err, false)
	}

	if err := verifyPassword(user.Password, data.Password); err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}

// UpdatePassword updates user's password by user ID
// if user exists and currentPassword is correct.
func (u *user) UpdatePassword(id uuid.UUID, currentPassword []byte, newPassword []byte) error {
	user, err := u.Get(id)
	if err != nil {
		return err
	}

	if err = verifyPassword(user.Password, currentPassword); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(newPassword, u.params.HashCost)
	if err != nil {
		return err
	}

	if err = u.repos.User.UpdatePassword(context.Background(), user.ID, hashedPassword); err != nil {
		return NewError(err, false)
	}

	return nil
}

// Delete deletes a user by ID
// if user exists and currentPassword is correct.
func (u *user) Delete(id uuid.UUID, currentPassword []byte) error {
	user, err := u.Get(id)
	if err != nil {
		return err
	}

	if err = verifyPassword(user.Password, currentPassword); err != nil {
		return err
	}

	if err := u.repos.User.DeleteByID(context.Background(), user.ID); err != nil {
		return NewError(err, false)
	}

	return nil
}
