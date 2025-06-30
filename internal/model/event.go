package model

import (
	"slices"
	"time"
)

const (
	eventDescLength = 32
)

type Event struct {
	date  time.Time
	money int32
	desc  string
	tags  []Tag
}

type EventWithID struct {
	Event
	id int32
}

func ParseDate(ds string) (*time.Time, error) {
	layouts := []string{
		"2006/1/2",
		"2006/01/02",
		"2006-01-02",
	}
	var err error
	for _, layout := range layouts {
		date, err := time.Parse(layout, ds)
		if err == nil {
			return &date, nil
		}
	}
	return nil, err
}

func NewEvent(date time.Time, money int32,
	desc string, tags []Tag) *Event {
	if len([]rune(desc)) >= eventDescLength {
		desc = string([]rune(desc)[0:eventDescLength])
	}
	return &Event{
		date:  date,
		money: money,
		desc:  desc,
		tags:  tags,
	}
}

func (e *Event) GetDate() time.Time {
	return e.date
}

func (e *Event) GetMoney() int32 {
	return e.money
}

func (e *Event) GetDesc() string {
	return e.desc
}

func (e *Event) GetTags() []Tag {
	ret := make([]Tag, len(e.tags))
	copy(ret, e.tags)
	return ret
}

func (e *Event) AddTag(tag Tag) {
	if !slices.Contains(e.tags, tag) {
		e.tags = append(e.tags, tag)
	}
}

func NewEventWithID(id int32, date time.Time, money int32,
	desc string, tags []Tag) *EventWithID {
	if len([]rune(desc)) >= eventDescLength {
		desc = string([]rune(desc)[0:eventDescLength])
	}
	return &EventWithID{
		id:    id,
		Event: *NewEvent(date, money, desc, tags),
	}
}

func (e *EventWithID) GetID() int32 {
	return e.id
}
