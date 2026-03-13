package shortener

import (
	"crypto/sha256"

	"github.com/jxskiss/base62"
)

// codeLength is the number of hash bytes to keep before base62 encoding.
// 6 bytes = 48 bits ≈ 281 trillion unique codes → ~8 char base62 string.
const codeLength = 6

// Encode produces a short, fixed-length base62 code from any input URL.
// It hashes the input with SHA-256, truncates to `codeLength` bytes, and
// base62-encodes the result.
func Encode(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base62.EncodeToString(hash[:codeLength])
}
