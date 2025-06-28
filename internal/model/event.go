package model

import (
	"slices"
	"time"
)

type Event struct {
	date  time.Time
	money int32
	desc  string
	tags  []Tag
}

func NewEvent(date time.Time, money int32,
	desc string, tags []Tag) *Event {
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

func (e *Event) AddTag(tags []Tag) {
	for _, tag := range tags {
		if !slices.Contains(e.tags, tag) {
			e.tags = append(e.tags, tag)
		}
	}
}
