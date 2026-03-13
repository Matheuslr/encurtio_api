package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/matheuslr/encurtio/internal/model"
)

type CassandraURLRepository struct {
	session *gocql.Session
}

func NewCassandraURLRepository(session *gocql.Session) *CassandraURLRepository {
	return &CassandraURLRepository{session: session}
}

func (r *CassandraURLRepository) Save(ctx context.Context, url *model.URL) error {
	query := `
			INSERT INTO urls (short_code, original_url, created_at)
			values(?,?,?)
		`

	return r.session.Query(
		query,
		url.ShortCode,
		url.Original,
		url.CreatedAt,
	).WithContext(ctx).Exec()
}

func (r *CassandraURLRepository) FindByShortCode(ctx context.Context, code string) (*model.URL, error) {
	var url model.URL

	query := `
		SELECT short_code, original_url, created_at
		FROM urls
		WHERE short_code = ? 
	`

	err := r.session.Query(query, code).
		WithContext(ctx).
		Scan(&url.ShortCode, &url.Original, &url.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil

}
