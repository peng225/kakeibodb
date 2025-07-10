package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/sqlc/query"
	"kakeibodb/internal/usecase"
	"slices"
	"time"
)

type EventRepository struct {
	q *query.Queries
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		q: query.New(db),
	}
}

func (er *EventRepository) Create(ctx context.Context, req *usecase.EventCreateRequest) (int64, error) {
	tx := txFromCtx(ctx)
	if tx == nil {
		return 0, errors.New("failed to get tx from context")
	}
	qtx := er.q.WithTx(tx)
	event, err := getByContent(ctx, qtx, req.Date, req.Money, req.Desc)
	if err == nil {
		return event.GetID(), nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	res, err := qtx.CreateEvent(ctx, query.CreateEventParams{
		Dt: sql.NullTime{
			Time:  req.Date,
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: req.Money,
			Valid: true,
		},
		Description: sql.NullString{
			String: req.Desc,
			Valid:  true,
		},
	})
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func getByContent(ctx context.Context, qtx *query.Queries,
	date time.Time, money int32, desc string) (*model.Event, error) {
	res, err := qtx.GetEvent(ctx, query.GetEventParams{
		Dt: sql.NullTime{
			Time:  date,
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: money,
			Valid: true,
		},
		Description: sql.NullString{
			String: desc,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return model.NewEvent(
		res.ID, res.Dt.Time, res.Money.Int32,
		res.Description.String, nil,
	), nil
}

func (er *EventRepository) GetWithoutTags(ctx context.Context, id int64) (*model.Event, error) {
	res, err := er.q.GetEventByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return model.NewEvent(res.ID, res.Dt.Time, res.Money.Int32, res.Description.String, nil), nil
}

func (er *EventRepository) UpdateMoney(ctx context.Context, id int64, money int32) error {
	tx := txFromCtx(ctx)
	if tx == nil {
		return errors.New("failed to get tx from context")
	}
	qtx := er.q.WithTx(tx)
	err := qtx.UpdateEventMoney(ctx, query.UpdateEventMoneyParams{
		ID: id,
		Money: sql.NullInt32{
			Int32: money,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update event's money: %w", err)
	}
	return nil
}

func (er *EventRepository) Delete(ctx context.Context, id int64) error {
	tx := txFromCtx(ctx)
	if tx == nil {
		return errors.New("failed to get tx from context")
	}
	qtx := er.q.WithTx(tx)
	err := qtx.DeleteEventByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete event by ID: %w", err)
	}
	return nil
}

func sortEventsByDate(events []*model.Event) {
	slices.SortFunc(events, func(a, b *model.Event) int {
		if a.GetDate().Before(b.GetDate()) {
			return -1
		} else if a.GetDate().After(b.GetDate()) {
			return 1
		}
		return 0
	})
}

func (er *EventRepository) ListOutcomes(ctx context.Context, from, to time.Time) ([]*model.Event, error) {
	toInclusive := to.AddDate(0, 0, -1)
	res, err := er.q.ListOutcomeEvents(ctx, query.ListOutcomeEventsParams{
		FromDt: sql.NullTime{
			Time:  from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  toInclusive,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list outcome events: %w", err)
	}

	events := make([]*model.Event, 0)
	for _, ewt := range res {
		e := model.NewEvent(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		if !ewt.Tagname.Valid {
			events = append(events, e)
			continue
		}
		tagName := ewt.Tagname.String
		if len(events) == 0 {
			e.AddTag(tagName)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetID() == lastEvent.GetID() {
				lastEvent.AddTag(tagName)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tagName)
				events = append(events, e)
			}
		}
	}
	sortEventsByDate(events)
	return events, nil
}

func (er *EventRepository) ListOutcomesWithTags(ctx context.Context, tagNames []string,
	from, to time.Time) ([]*model.Event, error) {
	sqlTags := make([]sql.NullString, len(tagNames))
	for i, tagName := range tagNames {
		sqlTags[i] = sql.NullString{
			String: tagName,
			Valid:  true,
		}
	}
	toInclusive := to.AddDate(0, 0, -1)
	res, err := er.q.ListOutcomeEventsWithTags(ctx, query.ListOutcomeEventsWithTagsParams{
		FromDt: sql.NullTime{
			Time:  from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  toInclusive,
			Valid: true,
		},
		Tagnames: sqlTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list outcome events with tags: %w", err)
	}

	events := make([]*model.Event, 0)
	for _, ewt := range res {
		e := model.NewEvent(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tagName := ewt.Tagname.String
		if len(events) == 0 {
			e.AddTag(tagName)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetID() == lastEvent.GetID() {
				lastEvent.AddTag(tagName)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tagName)
				events = append(events, e)
			}
		}
	}
	sortEventsByDate(events)
	return events, nil
}

func (er *EventRepository) List(ctx context.Context, from, to time.Time) ([]*model.Event, error) {
	toInclusive := to.AddDate(0, 0, -1)
	res, err := er.q.ListEvents(ctx, query.ListEventsParams{
		FromDt: sql.NullTime{
			Time:  from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  toInclusive,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	events := make([]*model.Event, 0)
	for _, ewt := range res {
		e := model.NewEvent(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tagName := ewt.Tagname.String
		if len(events) == 0 {
			e.AddTag(tagName)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetID() == lastEvent.GetID() {
				lastEvent.AddTag(tagName)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tagName)
				events = append(events, e)
			}
		}
	}
	sortEventsByDate(events)
	return events, nil
}

func (er *EventRepository) ListWithTags(ctx context.Context, tagNames []string,
	from, to time.Time) ([]*model.Event, error) {
	sqlTags := make([]sql.NullString, len(tagNames))
	for i, tagName := range tagNames {
		sqlTags[i] = sql.NullString{
			String: tagName,
			Valid:  true,
		}
	}
	toInclusive := to.AddDate(0, 0, -1)
	res, err := er.q.ListEventsWithTags(ctx, query.ListEventsWithTagsParams{
		FromDt: sql.NullTime{
			Time:  from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  toInclusive,
			Valid: true,
		},
		Tagnames: sqlTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events with tags: %w", err)
	}

	events := make([]*model.Event, 0)
	for _, ewt := range res {
		e := model.NewEvent(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tagName := ewt.Tagname.String
		if len(events) == 0 {
			e.AddTag(tagName)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetID() == lastEvent.GetID() {
				lastEvent.AddTag(tagName)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tagName)
				events = append(events, e)
			}
		}
	}
	sortEventsByDate(events)
	return events, nil
}
