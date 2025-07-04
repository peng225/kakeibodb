package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
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
	pwi, err := pr.getByKey(ctx, key)
	if err == nil {
		return pwi.GetID(), nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	res, err := pr.q.CreatePattern(ctx, sql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create pattern: %w", err)
	}
	return res.LastInsertId()
}

func (pr *PatternRepository) getByKey(ctx context.Context, key string) (*model.Pattern, error) {
	res, err := pr.q.GetPattern(ctx, sql.NullString{
		String: key,
		Valid:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pattern: %w", err)
	}
	return model.NewPattern(res.ID, res.KeyString.String, nil), nil
}

func (pr *PatternRepository) Delete(id int64) error {
	ctx := context.Background()
	err := pr.q.DeletePatternByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pattern by ID: %w", err)
	}
	return nil
}

func (pr *PatternRepository) List() ([]*model.Pattern, error) {
	ctx := context.Background()
	res, err := pr.q.ListPatterns(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list patterns: %w", err)
	}
	patterns := make([]*model.Pattern, 0)
	for _, pwt := range res {
		p := model.NewPattern(pwt.ID, pwt.KeyString.String, nil)
		tag := model.Tag(pwt.Tags.String)
		if len(patterns) == 0 {
			p.AddTag(tag)
			patterns = append(patterns, p)
		} else {
			lastPattern := patterns[len(patterns)-1]
			if p.GetKey() == lastPattern.GetKey() {
				lastPattern.AddTag(tag)
				patterns[len(patterns)-1] = lastPattern
			} else {
				p.AddTag(tag)
				patterns = append(patterns, p)
			}
		}
	}
	return patterns, nil
}
