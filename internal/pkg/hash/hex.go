package hash

import "encoding/hex"

type Hash []byte

func FromHexString(hexString string) Hash {
	hash, _ := hex.DecodeString(hexString)
	return hash
}

func (h Hash) HexString() string {
	return hex.EncodeToString(h)
}
