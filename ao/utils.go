package ao

import (
	"crypto/rand"
	"encoding/base64"
)

func Coalesce[T any](value interface{}, defaultValue any) T {
    if value != nil {
        return value.(T)
    }
    return defaultValue.(T)
}

// CreateRandomAnchor generates a random 32-byte string encoded as base64 URL-safe
func CreateRandomAnchor() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}