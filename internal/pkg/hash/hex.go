// Package hash provides structure to work with hex representation of hash.
package hash

import "encoding/hex"

// Hash represents byte slice with methods.
type Hash []byte

// FromHexString decodes hex string and returns Hash instance.
func FromHexString(hexString string) Hash {
	hash, _ := hex.DecodeString(hexString)
	return hash
}

// HexString encodes Hash instance and returns hex string.
func (h Hash) HexString() string {
	return hex.EncodeToString(h)
}
