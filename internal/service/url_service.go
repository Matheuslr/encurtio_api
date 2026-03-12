package service

import (
	"context"

	"github.com/matheuslr/encurtio/internal/model"
	"github.com/matheuslr/encurtio/internal/repository"
	"github.com/matheuslr/encurtio/internal/shortener"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Shorten(ctx context.Context, Original string) (*model.URL, error) {

	short_code := shortener.Encode(Original)
	url := &model.URL{
		ShortCode: short_code,
		Original:  Original,
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
