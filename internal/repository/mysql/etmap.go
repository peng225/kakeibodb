package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/repository/mysql/sqlc/query"
)

type EventTagMapRepository struct {
	q *query.Queries
}

func NewEventTagMapRepository(db *sql.DB) *EventTagMapRepository {
	return &EventTagMapRepository{
		q: query.New(db),
	}
}

func (etmr *EventTagMapRepository) Map(eventID int64, tagName string) error {
	ctx := context.Background()
	res, err := etmr.q.GetTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	tagID := res.ID

	_, err = etmr.q.GetEventToTagMap(ctx, query.GetEventToTagMapParams{
		EventID: sql.NullInt64{
			Int64: eventID,
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
		return fmt.Errorf("failed to get event-to-tag map: %w", err)
	}

	_, err = etmr.q.MapEventToTag(ctx, query.MapEventToTagParams{
		EventID: sql.NullInt64{
			Int64: eventID,
			Valid: true,
		},
		TagID: sql.NullInt64{
			Int64: tagID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to map event to tag: %w", err)
	}
	return nil
}

func (etmr *EventTagMapRepository) Unmap(eventID int64, tagName string) error {
	ctx := context.Background()
	res, err := etmr.q.GetTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	tagID := res.ID

	_, err = etmr.q.UnmapEventFromTag(ctx, query.UnmapEventFromTagParams{
		EventID: sql.NullInt64{
			Int64: eventID,
			Valid: true,
		},
		TagID: sql.NullInt64{
			Int64: tagID,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to unmap event from tag: %w", err)
	}
	return nil
}
