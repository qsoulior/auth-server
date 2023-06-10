package usecase

import (
	"bytes"
	"crypto/sha256"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

type fingerprint struct {
	data   []byte
	userID uuid.UUID
}

func NewFingerprint(data []byte, userID uuid.UUID) *fingerprint {
	return &fingerprint{data, userID}
}

func (f *fingerprint) Hash() ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(append(f.data, f.userID[:]...))
	if err != nil {
		return nil, ErrFingerprintInvalid
	}

	return h.Sum(nil), nil
}

func (f *fingerprint) Verify(hash []byte) error {
	h, err := f.Hash()
	if err != nil {
		return err
	}

	if !bytes.Equal(h, hash) {
		return ErrFingerprintIncorrect
	}
	return nil
}
