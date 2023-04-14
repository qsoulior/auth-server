package rand

import "crypto/rand"

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GetString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = chars[v%byte(len(chars))]
	}

	return string(b), err
}
