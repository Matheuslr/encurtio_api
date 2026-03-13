package shortener_test

import (
	"testing"

	"github.com/matheuslr/encurtio/internal/shortener"
	"github.com/stretchr/testify/assert"
)

func TestEncode_ProducesShortCode(t *testing.T) {
	urls := []string{
		"https://www.google.com/search?q=golang+url+shortener",
		"https://github.com/Matheuslr/encurtio_api",
		"https://example.com",
		"https://example.com/some/very/long/path?with=params&and=more&foo=bar",
	}

	for _, u := range urls {
		code := shortener.Encode(u)
		t.Logf("%-60s → %s (len=%d)", u, code, len(code))

		assert.NotEmpty(t, code, "code should not be empty for %q", u)
		assert.LessOrEqual(t, len(code), 10, "code should be ≤10 chars for %q", u)
	}
}

func TestEncode_IsDeterministic(t *testing.T) {
	input := "https://example.com/deterministic-test"
	a := shortener.Encode(input)
	b := shortener.Encode(input)
	assert.Equal(t, a, b, "same input must always produce the same code")
}

func TestEncode_DifferentInputsProduceDifferentCodes(t *testing.T) {
	a := shortener.Encode("https://example.com/page-a")
	b := shortener.Encode("https://example.com/page-b")
	assert.NotEqual(t, a, b, "different inputs should produce different codes")
}

func TestEncode_EmptyInput(t *testing.T) {
	code := shortener.Encode("")
	assert.NotEmpty(t, code, "even empty input should produce a code")
	assert.LessOrEqual(t, len(code), 10)
}

func TestEncode_AlwaysShorterThanInput(t *testing.T) {
	longURL := "https://example.com/very/long/path/that/goes/on/and/on?with=query&params=yes&more=stuff"
	code := shortener.Encode(longURL)
	assert.Less(t, len(code), len(longURL), "encoded code must be shorter than the original URL")
}
