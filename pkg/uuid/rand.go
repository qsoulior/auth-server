package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalid = errors.New("uuid: invalid UUID")
)

type UUID [16]byte

// V4, RFC4122
func New() (UUID, error) {
	var uuid UUID
	if _, err := rand.Read(uuid[:]); err != nil {
		return UUID{}, err
	}

	uuid[6] = uuid[6]&0x0f | 0x40
	uuid[8] = uuid[8]&0x3f | 0x80

	return uuid, nil
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UUID) Scan(src any) error {
	if s, ok := src.(string); ok {
		uuid, err := FromString(s)
		if err != nil {
			return ErrInvalid
		}
		*u = uuid
		return nil
	}
	return ErrInvalid
}

func FromString(s string) (UUID, error) {
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return UUID{}, ErrInvalid
	}

	s = strings.ReplaceAll(s, "-", "")

	var uuid UUID
	if _, err := hex.Decode(uuid[:], []byte(s)); err != nil {
		return UUID{}, ErrInvalid
	}
	return uuid, nil
}

func FromBytes(b []byte) (UUID, error) {
	if len(b) != 16 {
		return UUID{}, ErrInvalid
	}

	var uuid UUID
	copy(uuid[:], b)
	return uuid, nil
}
