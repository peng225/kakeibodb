package usecase

import (
	"kakeibodb/internal/model"
	"time"
)

type EventRepository interface {
	Create(req *EventCreateRequest) (int64, error)
	GetWithoutTags(id int64) (*model.Event, error)
	UpdateMoney(id int64, money int32) error
	Delete(id int64) error
	ListOutcomes(from, to time.Time) ([]*model.Event, error)
	ListOutcomesWithTags(tagNames []string, from, to time.Time) ([]*model.Event, error)
	List(from, to time.Time) ([]*model.Event, error)
	ListWithTags(tagNames []string, from, to time.Time) ([]*model.Event, error)
}

type EventTagMapRepository interface {
	Map(eventID int64, tagName string) error
	Unmap(eventID int64, tagName string) error
}

type TagRepository interface {
	Create(tag string) (int64, error)
	Delete(id int64) error
	List() ([]*model.Tag, error)
}

type PatternRepository interface {
	Create(key string) (int64, error)
	Delete(id int64) error
	List() ([]*model.Pattern, error)
}

type PatternTagMapRepository interface {
	Map(patternID int64, tagName string) error
	Unmap(patternID int64, tagName string) error
}
