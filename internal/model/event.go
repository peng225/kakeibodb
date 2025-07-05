package model

import (
	"slices"
	"time"
)

const (
	eventDescLength = 32
)

type Event struct {
	id       int64
	date     time.Time
	money    int32
	desc     string
	tagNames []string
}

func NewEvent(id int64, date time.Time, money int32,
	desc string, tagNames []string) *Event {
	desc = FormatDesc(desc)
	return &Event{
		id:       id,
		date:     date,
		money:    money,
		desc:     desc,
		tagNames: tagNames,
	}
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

func FormatDesc(desc string) string {
	if len([]rune(desc)) >= eventDescLength {
		desc = string([]rune(desc)[0:eventDescLength])
	}
	return desc
}

func (e *Event) GetID() int64 {
	return e.id
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

func (e *Event) GetTagNames() []string {
	ret := make([]string, len(e.tagNames))
	copy(ret, e.tagNames)
	return ret
}

func (e *Event) AddTag(tagName string) {
	if !slices.Contains(e.tagNames, tagName) {
		e.tagNames = append(e.tagNames, tagName)
	}
}
