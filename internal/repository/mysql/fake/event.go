package fake

import (
	"context"
	"errors"
	"kakeibodb/internal/model"
	"kakeibodb/internal/usecase"
	"slices"
	"time"
)

type EventFakeRepository struct {
	nextID int64
	events []*model.Event
}

func NewEventFakeRepository() *EventFakeRepository {
	return &EventFakeRepository{}
}

func (er *EventFakeRepository) Create(ctx context.Context, req *usecase.EventCreateRequest) (int64, error) {
	er.events = append(er.events, model.NewEvent(er.nextID, req.Date, req.Money, req.Desc, nil))
	er.nextID += 1
	return er.nextID - 1, nil
}

func (er *EventFakeRepository) GetWithoutTags(ctx context.Context, id int64) (*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) UpdateMoney(ctx context.Context, id int64, money int32) error {
	// Not implemented.
	return nil
}

func (er *EventFakeRepository) Delete(ctx context.Context, id int64) error {
	// Not implemented.
	return nil
}

func (er *EventFakeRepository) ListOutcomes(ctx context.Context, from, to time.Time) ([]*model.Event, error) {
	seq := func(yield func(*model.Event) bool) {
		for _, event := range er.events {
			if event.GetMoney() < 0 &&
				(event.GetDate().Equal(from) || event.GetDate().After(from)) &&
				event.GetDate().Before(to) && // "to" is exclusive edge.
				!yield(event) {
				return
			}
		}
	}

	matchedEvents := slices.Collect(seq)
	return matchedEvents, nil
}

func (er *EventFakeRepository) ListOutcomesWithTags(ctx context.Context, tagNames []string,
	from, to time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) List(ctx context.Context, from, to time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) ListWithTags(ctx context.Context, tagNames []string,
	from, to time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) AddTag(ctx context.Context, id int64, tagName string) error {
	if id >= er.nextID {
		return errors.New("invalid ID")
	}
	er.events[id].AddTag(tagName)
	return nil
}
