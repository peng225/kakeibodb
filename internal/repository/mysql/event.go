package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/query/query"
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

func (er *EventRepository) Create(event *model.Event) (int64, error) {
	ctx := context.Background()
	res, err := er.q.CreateEvent(ctx, query.CreateEventParams{
		Dt: sql.NullTime{
			Time:  event.GetDate(),
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: event.GetMoney(),
			Valid: true,
		},
		Description: sql.NullString{
			String: event.GetDesc(),
			Valid:  true,
		},
	})
	if err != nil {
		return 0, err
	}
	// FIXME: tags are ignored. Is the model correct?
	return res.LastInsertId()
}

func (er *EventRepository) Exist(event *model.Event) (bool, error) {
	ctx := context.Background()
	_, err := er.q.GetEvent(ctx, query.GetEventParams{
		Dt: sql.NullTime{
			Time:  event.GetDate(),
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: event.GetMoney(),
			Valid: true,
		},
		Description: sql.NullString{
			String: event.GetDesc(),
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (er *EventRepository) Get(id int32) (*model.Event, error) {
	ctx := context.Background()
	res, err := er.q.GetEventByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return model.NewEvent(res.Dt.Time, res.Money.Int32, res.Description.String, nil), nil
}

func (er *EventRepository) Delete(id int32) error {
	ctx := context.Background()
	err := er.q.DeleteEventByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete event by ID: %w", err)
	}
	return nil
}

func (er *EventRepository) ListOutcomes(from, to *time.Time) ([]*model.EventWithID, error) {
	ctx := context.Background()
	res, err := er.q.ListOutcomeEvents(ctx, query.ListOutcomeEventsParams{
		FromDt: sql.NullTime{
			Time:  *from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  *to,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list outcome events: %w", err)
	}

	events := make([]*model.EventWithID, 0)
	for _, ewt := range res {
		e := model.NewEventWithID(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tag := model.Tag(ewt.Tags.String)
		if len(events) == 0 {
			e.AddTag(tag)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetDate().Equal(lastEvent.GetDate()) &&
				e.GetMoney() == lastEvent.GetMoney() &&
				e.GetDesc() == lastEvent.GetDesc() {
				lastEvent.AddTag(tag)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tag)
				events = append(events, e)
			}
		}
	}
	return events, nil
}

func (er *EventRepository) ListOutcomesWithTags(tags []model.Tag, from, to *time.Time) ([]*model.EventWithID, error) {
	sqlTags := make([]sql.NullString, len(tags))
	for i, tag := range tags {
		sqlTags[i] = sql.NullString{
			String: tag.String(),
			Valid:  true,
		}
	}
	ctx := context.Background()
	res, err := er.q.ListOutcomeEventsWithTags(ctx, query.ListOutcomeEventsWithTagsParams{
		FromDt: sql.NullTime{
			Time:  *from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  *to,
			Valid: true,
		},
		Tags: sqlTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list outcome events with tags: %w", err)
	}

	events := make([]*model.EventWithID, 0)
	for _, ewt := range res {
		e := model.NewEventWithID(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tag := model.Tag(ewt.Tags.String)
		if len(events) == 0 {
			e.AddTag(tag)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetDate().Equal(lastEvent.GetDate()) &&
				e.GetMoney() == lastEvent.GetMoney() &&
				e.GetDesc() == lastEvent.GetDesc() {
				lastEvent.AddTag(tag)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tag)
				events = append(events, e)
			}
		}
	}
	return events, nil
}

func (er *EventRepository) List(from, to *time.Time) ([]*model.EventWithID, error) {
	ctx := context.Background()
	res, err := er.q.ListEvents(ctx, query.ListEventsParams{
		FromDt: sql.NullTime{
			Time:  *from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  *to,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	events := make([]*model.EventWithID, 0)
	for _, ewt := range res {
		e := model.NewEventWithID(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tag := model.Tag(ewt.Tags.String)
		if len(events) == 0 {
			e.AddTag(tag)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetDate().Equal(lastEvent.GetDate()) &&
				e.GetMoney() == lastEvent.GetMoney() &&
				e.GetDesc() == lastEvent.GetDesc() {
				lastEvent.AddTag(tag)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tag)
				events = append(events, e)
			}
		}
	}
	return events, nil
}

func (er *EventRepository) ListWithTags(tags []model.Tag, from, to *time.Time) ([]*model.EventWithID, error) {
	sqlTags := make([]sql.NullString, len(tags))
	for i, tag := range tags {
		sqlTags[i] = sql.NullString{
			String: tag.String(),
			Valid:  true,
		}
	}
	ctx := context.Background()
	res, err := er.q.ListEventsWithTags(ctx, query.ListEventsWithTagsParams{
		FromDt: sql.NullTime{
			Time:  *from,
			Valid: true,
		},
		ToDt: sql.NullTime{
			Time:  *to,
			Valid: true,
		},
		Tags: sqlTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events with tags: %w", err)
	}

	events := make([]*model.EventWithID, 0)
	for _, ewt := range res {
		e := model.NewEventWithID(ewt.ID, ewt.Dt.Time, ewt.Money.Int32, ewt.Description.String, nil)
		tag := model.Tag(ewt.Tags.String)
		if len(events) == 0 {
			e.AddTag(tag)
			events = append(events, e)
		} else {
			lastEvent := events[len(events)-1]
			if e.GetDate().Equal(lastEvent.GetDate()) &&
				e.GetMoney() == lastEvent.GetMoney() &&
				e.GetDesc() == lastEvent.GetDesc() {
				lastEvent.AddTag(tag)
				events[len(events)-1] = lastEvent
			} else {
				e.AddTag(tag)
				events = append(events, e)
			}
		}
	}
	return events, nil
}
