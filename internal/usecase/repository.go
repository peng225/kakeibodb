package usecase

import (
	"context"
	"kakeibodb/internal/model"
	"time"
)

type EventRepository interface {
	Create(ctx context.Context, req *EventCreateRequest) (int64, error)
	GetWithoutTags(ctx context.Context, id int64) (*model.Event, error)
	UpdateMoney(ctx context.Context, id int64, money int32) error
	Delete(ctx context.Context, id int64) error
	ListOutcomes(ctx context.Context, from, to time.Time) ([]*model.Event, error)
	ListOutcomesWithTags(ctx context.Context, tagNames []string, from, to time.Time) ([]*model.Event, error)
	List(ctx context.Context, from, to time.Time) ([]*model.Event, error)
	ListWithTags(ctx context.Context, tagNames []string, from, to time.Time) ([]*model.Event, error)
}

type EventTagMapRepository interface {
	Map(ctx context.Context, eventID int64, tagName string) error
	Unmap(ctx context.Context, eventID int64, tagName string) error
}

type TagRepository interface {
	Create(ctx context.Context, tag string) (int64, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*model.Tag, error)
}

type PatternRepository interface {
	Create(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*model.Pattern, error)
}

type PatternTagMapRepository interface {
	Map(ctx context.Context, patternID int64, tagName string) error
	Unmap(ctx context.Context, patternID int64, tagName string) error
}
