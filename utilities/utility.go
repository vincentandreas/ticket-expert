package utilities

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashParams(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}
