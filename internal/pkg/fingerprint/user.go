// Package fingerprint provides structure to hash and verify fingerprint.
package fingerprint

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/qsoulior/auth-server/internal/pkg/hash"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

var (
	ErrFingerprintIncorrect = errors.New("fingerprint is incorrect")
	ErrFingerprintInvalid   = errors.New("fingerprint is invalid")
)

// Fingerprint structure.
type Fingerprint struct {
	id   uuid.UUID
	data []byte
}

// New returns pointer to a Fingerprint instance.
func New(id uuid.UUID, data []byte) *Fingerprint {
	return &Fingerprint{id, data}
}

// Hash returns hash of fingerprint with salt.
func (f *Fingerprint) Hash() (hash.Hash, error) {
	h := sha256.New()
	_, err := h.Write(append(f.id[:], f.data...))
	if err != nil {
		return nil, ErrFingerprintInvalid
	}

	return h.Sum(nil), nil
}

// Verify compares fingerprint hash with a hash argument.
// It returns nil if hashes are equal.
func (f *Fingerprint) Verify(hash hash.Hash) error {
	h, err := f.Hash()
	if err != nil {
		return err
	}

	if !bytes.Equal(h, hash) {
		return ErrFingerprintIncorrect
	}
	return nil
}
