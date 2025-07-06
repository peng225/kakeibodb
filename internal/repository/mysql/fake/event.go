package fake

import (
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

func (er *EventFakeRepository) Create(req *usecase.EventCreateRequest) (int64, error) {
	er.events = append(er.events, model.NewEvent(er.nextID, req.Date, req.Money, req.Desc, nil))
	er.nextID += 1
	return er.nextID - 1, nil
}

func (er *EventFakeRepository) GetWithoutTags(id int64) (*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) UpdateMoney(id int64, money int32) error {
	// Not implemented.
	return nil
}

func (er *EventFakeRepository) Delete(id int64) error {
	// Not implemented.
	return nil
}

func (er *EventFakeRepository) ListOutcomes(from, to *time.Time) ([]*model.Event, error) {
	seq := func(yield func(*model.Event) bool) {
		for _, event := range er.events {
			if event.GetMoney() < 0 &&
				(event.GetDate().Equal(*from) || event.GetDate().After(*from)) &&
				(event.GetDate().Equal(*to) || event.GetDate().Before(*to)) &&
				!yield(event) {
				return
			}
		}
	}

	matchedEvents := slices.Collect(seq)
	return matchedEvents, nil
}

func (er *EventFakeRepository) ListOutcomesWithTags(tagNames []string, from, to *time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) List(from, to *time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) ListWithTags(tagNames []string, from, to *time.Time) ([]*model.Event, error) {
	// Not implemented.
	return nil, nil
}

func (er *EventFakeRepository) AddTag(id int64, tagName string) error {
	if id >= er.nextID {
		return errors.New("invalid ID")
	}
	er.events[id].AddTag(tagName)
	return nil
}
