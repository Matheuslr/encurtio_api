package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matheuslr/encurtio/configs"
	"github.com/matheuslr/encurtio/internal/model"
	"github.com/matheuslr/encurtio/internal/service"
	"github.com/matheuslr/encurtio/internal/shortener"
	"github.com/matheuslr/encurtio/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestService() (*service.URLService, *mocks.MockURLRepository) {
	repo := new(mocks.MockURLRepository)
	cfg := configs.Config{
		API: configs.APIConfig{
			Port: "8080",
			URL:  "http://localhost:8080",
		},
	}
	svc := service.NewURLService(repo, cfg)
	return svc, repo
}

// --- Shorten ---

func TestShorten_Success(t *testing.T) {
	svc, repo := newTestService()

	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.URL")).Return(nil)

	url, err := svc.Shorten(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, "https://example.com", url.Original)
	assert.Equal(t, shortener.Encode("https://example.com"), url.ShortCode)
	assert.WithinDuration(t, time.Now(), url.CreatedAt, 2*time.Second)

	repo.AssertExpectations(t)
}

func TestShorten_RepoError(t *testing.T) {
	svc, repo := newTestService()

	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.URL")).
		Return(errors.New("db connection lost"))

	url, err := svc.Shorten(context.Background(), "https://example.com")

	assert.Error(t, err)
	assert.Nil(t, url)
	assert.EqualError(t, err, "db connection lost")

	repo.AssertExpectations(t)
}

// --- GetOriginalURL ---

func TestGetOriginalURL_Found(t *testing.T) {
	svc, repo := newTestService()

	expected := &model.URL{
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}
	repo.On("FindByShortCode", mock.Anything, "abc123").Return(expected, nil)

	url, err := svc.GetOriginalURL(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, expected, url)

	repo.AssertExpectations(t)
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	svc, repo := newTestService()

	repo.On("FindByShortCode", mock.Anything, "notfound").
		Return(nil, errors.New("not found"))

	url, err := svc.GetOriginalURL(context.Background(), "notfound")

	assert.Error(t, err)
	assert.Nil(t, url)

	repo.AssertExpectations(t)
}

// --- BuildShortURL ---

func TestBuildShortURL(t *testing.T) {
	svc, _ := newTestService()

	result := svc.BuildShortURL("abc123")
	assert.Equal(t, "http://localhost:8080/abc123", result)
}

func TestBuildShortURL_EmptyCode(t *testing.T) {
	svc, _ := newTestService()

	result := svc.BuildShortURL("")
	assert.Equal(t, "http://localhost:8080/", result)
}
