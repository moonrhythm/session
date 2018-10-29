package session

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func sign(value string, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write([]byte(value))
	digest := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(digest)
}

func verify(value, digest string, keys [][]byte) bool {
	for _, k := range keys {
		if hmac.Equal([]byte(digest), []byte(sign(value, k))) {
			return true
		}
	}
	return false
}
