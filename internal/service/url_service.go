package service

import (
	"context"
	"time"

	"github.com/matheuslr/encurtio/configs"
	"github.com/matheuslr/encurtio/internal/model"
	"github.com/matheuslr/encurtio/internal/repository"
	"github.com/matheuslr/encurtio/internal/shortener"
)

type URLService struct {
	repo   repository.URLRepository
	config configs.Config
}

func NewURLService(repo repository.URLRepository, config configs.Config) *URLService {
	return &URLService{repo: repo, config: config}
}

func (s *URLService) Shorten(ctx context.Context, Original string) (*model.URL, error) {

	short_code := shortener.Encode(Original)
	url := &model.URL{
		ShortCode: short_code,
		Original:  Original,
		CreatedAt: time.Now(),
	}

	err := s.repo.Save(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, code string) (*model.URL, error) {
	return s.repo.FindByShortCode(ctx, code)
}

func (s *URLService) BuildShortURL(shortCode string) string {
	return s.config.API.URL + "/" + shortCode
}
