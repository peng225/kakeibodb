package usecase

import (
	"fmt"
	"kakeibodb/internal/event"
	"kakeibodb/internal/model"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type EventCreateRequest struct {
	Date  time.Time
	Money int32
	Desc  string
}

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

type EventPresenter interface {
	Present(events []*model.Event)
}

type EventUseCase struct {
	eventRepo EventRepository
}

type EventPresentUseCase struct {
	EventUseCase
	eventPresenter EventPresenter
}

type EventTagMapUsecase struct {
	etmRepo EventTagMapRepository
}

type ApplyPatternUseCase struct {
	EventUseCase
	EventTagMapUsecase
	patternRepo PatternRepository
}

func NewEventUseCase(eventRepo EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
	}
}

func NewEventPresentUseCase(eventRepo EventRepository, eventPresenter EventPresenter) *EventPresentUseCase {
	return &EventPresentUseCase{
		EventUseCase:   *NewEventUseCase(eventRepo),
		eventPresenter: eventPresenter,
	}
}

func NewEventTagMapUseCase(etmRepo EventTagMapRepository) *EventTagMapUsecase {
	return &EventTagMapUsecase{
		etmRepo: etmRepo,
	}
}

func NewApplyPatternUseCase(eventRepo EventRepository, etmRepo EventTagMapRepository,
	patternRepo PatternRepository) *ApplyPatternUseCase {
	return &ApplyPatternUseCase{
		EventUseCase:       *NewEventUseCase(eventRepo),
		EventTagMapUsecase: *NewEventTagMapUseCase(etmRepo),
		patternRepo:        patternRepo,
	}
}

func (eu *EventUseCase) LoadFromFile(file string) error {
	// FIXME: Don't want to depend on a specific file format.
	csv := event.NewCSV()
	csv.Open(file)

	slog.Info("Load events from file.", "file", file)

	// Skip header
	_ = csv.Read()
	for {
		event := csv.Read()
		if event == nil {
			break
		}
		date, err := model.ParseDate(event[0])
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
		decrease := event[1]
		increase := event[2]
		desc := event[3]

		var money int32
		if (decrease == "" && increase == "") || (decrease != "" && increase != "") {
			return fmt.Errorf("bad event record. decrease = %s, increase = %s", decrease, increase)
		} else if decrease != "" {
			tmpMoney, err := strconv.ParseInt(decrease, 10, 32)
			if err != nil {
				return fmt.Errorf(`failed to parse "%s" as int: %w`, decrease, err)
			}
			money = -1 * int32(tmpMoney)
		} else {
			tmpMoney, err := strconv.ParseInt(increase, 10, 32)
			if err != nil {
				return fmt.Errorf(`failed to parse "%s" as int: %w`, increase, err)
			}
			money = int32(tmpMoney)
		}
		slog.Info("Create value.", "date", date,
			"money", money, "desc", desc)
		_, err = eu.eventRepo.Create(&EventCreateRequest{
			Date:  *date,
			Money: money,
			Desc:  model.FormatDesc(desc),
		})
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}
	}
	return nil
}

func (eu *EventUseCase) LoadFromDir(dir string) error {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk dir: %w", err)
	}

	for _, file := range files {
		err = eu.LoadFromFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (eu *EventUseCase) LoadCreditFromFile(file string, relatedEventID int64) error {
	csv := event.NewCSV()
	csv.Open(file)

	slog.Info("Load credit events from file.", "file", file)

	relatedEvent, err := eu.eventRepo.GetWithoutTags(relatedEventID)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Skip header
	_ = csv.Read()
	creditEventCreateReqs := make([]*EventCreateRequest, 0)
	for {
		event := csv.Read()
		if event == nil {
			break
		}

		if event[0] == "" {
			continue
		}
		date, err := model.ParseDate(event[0])
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
		desc := event[1]
		tmpMoney, err := strconv.ParseInt(event[2], 10, 32)
		if err != nil {
			return fmt.Errorf(`failed to parse "%s" as int: %w`, event[2], err)
		}
		tmpMoney *= -1
		money := int32(tmpMoney)

		creditEventCreateReqs = append(creditEventCreateReqs, &EventCreateRequest{
			Date:  *date,
			Money: money,
			Desc:  model.FormatDesc(desc),
		})
	}

	if !eu.deletingCorrectEvent(relatedEvent.GetMoney(), creditEventCreateReqs) {
		return fmt.Errorf("sum of credit events does not match original event (ID=%d)", relatedEventID)
	}
	for _, cecReq := range creditEventCreateReqs {
		slog.Info("Create value.", "date", cecReq.Date,
			"money", cecReq.Money, "desc", cecReq.Desc)
		_, err = eu.eventRepo.Create(cecReq)
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}
	}
	err = eu.eventRepo.Delete(relatedEventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (eu *EventUseCase) deletingCorrectEvent(relatedEventMoney int32, creditEventCreateReqs []*EventCreateRequest) bool {
	moneySum := int32(0)
	for _, cecReq := range creditEventCreateReqs {
		moneySum += cecReq.Money
	}

	return moneySum == relatedEventMoney
}

func (eu *EventUseCase) getEventIDFromSplitBaseTag(splitBaseTagName string,
	date time.Time) (int64, error) {
	from := date.AddDate(0, -2, -5)
	events, err := eu.eventRepo.ListOutcomesWithTags([]string{splitBaseTagName}, from, date)
	if err != nil {
		return 0, err
	}
	return events[len(events)-1].GetID(), nil
}

func (eu *EventUseCase) Split(eventID int64, splitBaseTagName string, date time.Time, money int32, desc string) error {
	if eventID == -1 {
		var err error
		eventID, err = eu.getEventIDFromSplitBaseTag(splitBaseTagName, date)
		if err != nil {
			return err
		}
		slog.Info("Auto detected.", "eventID", eventID)
	}

	event, err := eu.eventRepo.GetWithoutTags(eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event.GetMoney()*money <= 0 {
		return fmt.Errorf("Income/Outcome event can be split only by another income/outcome event.")
	}
	if math.Abs(float64(event.GetMoney())) < math.Abs(float64(money)) {
		return fmt.Errorf("abs(to be split event money)(%f) should be larger than or equal to abs(splitting event money)(%f)",
			math.Abs(float64(event.GetMoney())), math.Abs(float64(money)))
	}

	// Update the existing event.
	if event.GetMoney() == money {
		err = eu.eventRepo.Delete(eventID)
		if err != nil {
			return fmt.Errorf("failed to delete event: %w", err)
		}
	} else {
		err = eu.eventRepo.UpdateMoney(eventID, event.GetMoney()-money)
		if err != nil {
			return fmt.Errorf("failed to update event money: %w", err)
		}
	}

	// Insert a new event.
	_, err = eu.eventRepo.Create(&EventCreateRequest{
		Date:  date,
		Money: money,
		Desc:  model.FormatDesc(desc),
	})
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (eu *EventPresentUseCase) PresentOutcomes(tagNames []string, from, to time.Time) error {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.ListOutcomes(from, to)
		if err != nil {
			return fmt.Errorf("failed to list outcome events: %w", err)
		}
	} else {
		events, err = eu.eventRepo.ListOutcomesWithTags(tagNames, from, to)
		if err != nil {
			return fmt.Errorf("failed to list outcome events with tags: %w", err)
		}
	}

	eu.eventPresenter.Present(events)
	return nil
}

func (eu *EventPresentUseCase) PresentAll(tagNames []string, from, to time.Time) error {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.List(from, to)
		if err != nil {
			return fmt.Errorf("failed to list all events: %w", err)
		}
	} else {
		events, err = eu.eventRepo.ListWithTags(tagNames, from, to)
		if err != nil {
			return fmt.Errorf("failed to list all events with tags: %w", err)
		}
	}

	eu.eventPresenter.Present(events)
	return nil
}

func (etmu *EventTagMapUsecase) AddTag(eventID int64, tagNames []string) error {
	for _, tagName := range tagNames {
		err := etmu.etmRepo.Map(eventID, tagName)
		if err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}
	}
	return nil
}

func (etmu *EventTagMapUsecase) RemoveTag(eventID int64, tagName string) error {
	err := etmu.etmRepo.Unmap(eventID, tagName)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}
	return nil
}

func (apu *ApplyPatternUseCase) ApplyPattern(from, to time.Time) error {
	events, err := apu.eventRepo.List(from, to)
	if err != nil {
		return err
	}

	patterns, err := apu.patternRepo.List()
	if err != nil {
		return err
	}

	for _, event := range events {
		for _, pattern := range patterns {
			if strings.Contains(event.GetDesc(), pattern.GetKey()) {
				for _, tagName := range pattern.GetTagNames() {
					err = apu.etmRepo.Map(event.GetID(), tagName)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
