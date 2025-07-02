package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/repository/mysql/sqlc/query"
)

type PatternTagMapRepository struct {
	q *query.Queries
}

func NewPatternTagMapRepository(db *sql.DB) *PatternTagMapRepository {
	return &PatternTagMapRepository{
		q: query.New(db),
	}
}

func (ptmr *PatternTagMapRepository) Map(patternID int64, tagName string) error {
	ctx := context.Background()
	res, err := ptmr.q.GetTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	tagID := res.ID

	_, err = ptmr.q.GetPatternToTagMap(ctx, query.GetPatternToTagMapParams{
		PatternID: sql.NullInt64{
			Int64: patternID,
			Valid: true,
		},
		TagID: sql.NullInt64{
			Int64: tagID,
			Valid: true,
		},
	})
	if err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get pattern-to-tag map: %w", err)
	}

	_, err = ptmr.q.MapPatternToTag(ctx, query.MapPatternToTagParams{
		PatternID: sql.NullInt64{
			Int64: patternID,
			Valid: true,
		},
		TagID: sql.NullInt64{
			Int64: tagID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to map pattern to tag: %w", err)
	}
	return nil
}

func (ptmr *PatternTagMapRepository) Unmap(patternID int64, tagName string) error {
	ctx := context.Background()
	res, err := ptmr.q.GetTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	tagID := res.ID

	_, err = ptmr.q.UnmapPatternFromTag(ctx, query.UnmapPatternFromTagParams{
		PatternID: sql.NullInt64{
			Int64: patternID,
			Valid: true,
		},
		TagID: sql.NullInt64{
			Int64: tagID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to unmap pattern from tag: %w", err)
	}
	return nil
}
