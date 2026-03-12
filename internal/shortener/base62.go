package shortener

import "github.com/jxskiss/base62"

func Encode(input string) string {
	return base62.EncodeToString([]byte(input))
}
