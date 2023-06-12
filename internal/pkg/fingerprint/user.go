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

type Fingerprint struct {
	id   uuid.UUID
	data []byte
}

func New(id uuid.UUID, data []byte) *Fingerprint {
	return &Fingerprint{id, data}
}

func (f *Fingerprint) Hash() (hash.Hash, error) {
	h := sha256.New()
	_, err := h.Write(append(f.id[:], f.data...))
	if err != nil {
		return nil, ErrFingerprintInvalid
	}

	return h.Sum(nil), nil
}

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
