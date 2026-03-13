package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheuslr/encurtio/configs"
	"github.com/matheuslr/encurtio/internal/handler"
	"github.com/matheuslr/encurtio/internal/model"
	"github.com/matheuslr/encurtio/internal/service"
	"github.com/matheuslr/encurtio/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(repo *mocks.MockURLRepository) *gin.Engine {
	cfg := configs.Config{
		API: configs.APIConfig{Port: "8080", URL: "http://localhost:8080"},
	}
	svc := service.NewURLService(repo, cfg)
	h := handler.NewURLHandler(svc)

	r := gin.New()
	r.POST("/api/v1/url/shorten", h.Shorten)
	r.GET("/:code", h.GetRedirectURL)
	return r
}

// --- Shorten handler ---

func TestShorten_Success(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.URL")).Return(nil)

	router := setupRouter(repo)

	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "short_url")
	assert.Contains(t, resp["short_url"], "http://localhost:8080/")

	repo.AssertExpectations(t)
}

func TestShorten_InvalidBody(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/shorten", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", resp["error"])
}

func TestShorten_EmptyBody(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/shorten", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShorten_ServiceError(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.URL")).
		Return(errors.New("db error"))

	router := setupRouter(repo)

	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}

// --- GetRedirectURL handler ---

func TestGetRedirectURL_Success(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	expected := &model.URL{
		ShortCode: "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}
	repo.On("FindByShortCode", mock.Anything, "abc123").Return(expected, nil)

	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))

	repo.AssertExpectations(t)
}

func TestGetRedirectURL_NotFound(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	repo.On("FindByShortCode", mock.Anything, "notfound").
		Return(nil, errors.New("not found"))

	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "short code not found", resp["error"])

	repo.AssertExpectations(t)
}

// --- Context cancellation ---

func TestShorten_ContextCancelled(t *testing.T) {
	repo := new(mocks.MockURLRepository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.URL")).
		Return(context.Canceled)

	router := setupRouter(repo)

	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	repo.AssertExpectations(t)
}
