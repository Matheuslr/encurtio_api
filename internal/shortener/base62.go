package shortener

import (
	"crypto/sha256"

	"github.com/jxskiss/base62"
)

const codeLength = 6

func Encode(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base62.EncodeToString(hash[:codeLength])
}
