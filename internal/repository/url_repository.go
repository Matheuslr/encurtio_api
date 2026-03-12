package repository

import (
	"context"

	"github.com/matheuslr/encurtio/internal/model"
)

type URLRepository interface {
	Save(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, code string) (*model.URL, error)
}
