package fake

import (
	"context"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/usecase"
	"slices"
	"time"
)

type EventForTest struct {
	date     time.Time
	money    int32
	desc     string
	tagNames []string
}

func NewEventForTest(date time.Time, money int32, desc string) *EventForTest {
	return &EventForTest{
		date:     date,
		money:    money,
		desc:     desc,
		tagNames: make([]string, 0),
	}
}

type EventFakeRepository struct {
	nextID int64
	events map[int64]*EventForTest
}

func NewEventFakeRepository() *EventFakeRepository {
	return &EventFakeRepository{
		nextID: 0,
		events: make(map[int64]*EventForTest),
	}
}

func (er *EventFakeRepository) Create(ctx context.Context, req *usecase.EventCreateRequest) (int64, error) {
	er.events[er.nextID] = NewEventForTest(req.Date, req.Money, req.Desc)
	er.nextID += 1
	return er.nextID - 1, nil
}

func (er *EventFakeRepository) GetWithoutTags(ctx context.Context, id int64) (*model.Event, error) {
	event, ok := er.events[id]
	if !ok {
		return nil, fmt.Errorf("event with ID %d not found", id)
	}
	ret := model.NewEvent(id, event.date, event.money, event.desc, event.tagNames)
	return ret, nil
}

func (er *EventFakeRepository) UpdateMoney(ctx context.Context, id int64, money int32) error {
	er.events[id].money = money
	return nil
}

func (er *EventFakeRepository) Delete(ctx context.Context, id int64) error {
	delete(er.events, id)
	return nil
}

func (er *EventFakeRepository) ListOutcomes(ctx context.Context, from, to time.Time) ([]*model.Event, error) {
	matchedEvents := make(map[int64]*EventForTest)
	for id, event := range er.events {
		if event.money < 0 &&
			(event.date.Equal(from) || event.date.After(from)) &&
			event.date.Before(to) {
			matchedEvents[id] = event
		}
	}

	ret := make([]*model.Event, 0)
	for id, event := range matchedEvents {
		ret = append(ret, model.NewEvent(id, event.date, event.money, event.desc, event.tagNames))
	}
	slices.SortStableFunc(ret, func(a, b *model.Event) int {
		if a.GetDate().Before(b.GetDate()) {
			return -1
		} else if a.GetDate().After(b.GetDate()) {
			return 1
		}
		return 0
	})
	return ret, nil
}

func haveCommonTagName(tags0, tags1 []string) bool {
	for _, t0 := range tags0 {
		for _, t1 := range tags1 {
			if t0 == t1 {
				return true
			}
		}
	}
	return false
}

func (er *EventFakeRepository) ListOutcomesWithTags(ctx context.Context, tagNames []string,
	from, to time.Time) ([]*model.Event, error) {
	matchedEvents := make(map[int64]*EventForTest)
	for id, event := range er.events {
		if event.money < 0 &&
			(event.date.Equal(from) || event.date.After(from)) &&
			event.date.Before(to) {
			if haveCommonTagName(event.tagNames, tagNames) {
				matchedEvents[id] = event
			}
		}
	}

	ret := make([]*model.Event, 0)
	for id, event := range matchedEvents {
		ret = append(ret, model.NewEvent(id, event.date, event.money, event.desc, event.tagNames))
	}
	slices.SortStableFunc(ret, func(a, b *model.Event) int {
		if a.GetDate().Before(b.GetDate()) {
			return -1
		} else if a.GetDate().After(b.GetDate()) {
			return 1
		}
		return 0
	})
	return ret, nil
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
	er.events[id].tagNames = append(er.events[id].tagNames, tagName)
	return nil
}
