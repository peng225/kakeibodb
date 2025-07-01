package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/repository/mysql/sqlc/query"
)

type PatternRepository struct {
	q *query.Queries
}

func NewPatternRepository(db *sql.DB) *PatternRepository {
	return &PatternRepository{
		q: query.New(db),
	}
}

func (pr *PatternRepository) Create(key string) (int64, error) {
	ctx := context.Background()
	res, err := pr.q.CreatePattern(ctx, sql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get pattern: %w", err)
	}
	return res.LastInsertId()
}

func (pr *PatternRepository) Exist(key string) (bool, error) {
	ctx := context.Background()
	_, err := pr.q.GetPattern(ctx, sql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get pattern: %w", err)
	}
	return true, nil
}

func (pr *PatternRepository) Delete(id int64) error {
	ctx := context.Background()
	err := pr.q.DeletePatternByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pattern by ID: %w", err)
	}
	return nil
}
