package usecase

import (
	"context"
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

type EventUseCase struct {
	eventRepo EventRepository
	tx        Transaction
}

type EventPresentUseCase struct {
	eventRepo      EventRepository
	eventPresenter EventPresenter
}

type EventTagMapUsecase struct {
	etmRepo EventTagMapRepository
}

type ApplyPatternUseCase struct {
	EventUseCase
	etmRepo     EventTagMapRepository
	patternRepo PatternRepository
}

func NewEventUseCase(eventRepo EventRepository, tx Transaction) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
		tx:        tx,
	}
}

func NewEventPresentUseCase(eventRepo EventRepository, eventPresenter EventPresenter) *EventPresentUseCase {
	return &EventPresentUseCase{
		eventRepo:      eventRepo,
		eventPresenter: eventPresenter,
	}
}

func NewEventTagMapUseCase(etmRepo EventTagMapRepository) *EventTagMapUsecase {
	return &EventTagMapUsecase{
		etmRepo: etmRepo,
	}
}

func NewApplyPatternUseCase(eventRepo EventRepository, etmRepo EventTagMapRepository,
	patternRepo PatternRepository, tx Transaction) *ApplyPatternUseCase {
	return &ApplyPatternUseCase{
		EventUseCase: *NewEventUseCase(eventRepo, tx),
		etmRepo:      etmRepo,
		patternRepo:  patternRepo,
	}
}

func (eu *EventUseCase) LoadFromFile(ctx context.Context, file string) error {
	// FIXME: Don't want to depend on a specific file format.
	csv := event.NewCSV()
	csv.Open(file)

	slog.Info("Load events from file.", "file", file)

	// Skip header
	_ = csv.Read()
	err := eu.tx.Do(ctx, func(ctx context.Context) error {
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
			_, err = eu.eventRepo.Create(ctx, &EventCreateRequest{
				Date:  *date,
				Money: money,
				Desc:  model.FormatDesc(desc),
			})
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}
		}
		return nil
	})

	return err
}

func (eu *EventUseCase) LoadFromDir(ctx context.Context, dir string) error {
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
		err = eu.LoadFromFile(ctx, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (eu *EventUseCase) LoadCreditFromFile(ctx context.Context, file string, relatedEventID int64) error {
	csv := event.NewCSV()
	csv.Open(file)

	slog.Info("Load credit events from file.", "file", file)

	err := eu.tx.Do(ctx, func(ctx context.Context) error {
		relatedEvent, err := eu.eventRepo.GetWithoutTags(ctx, relatedEventID)
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
			_, err = eu.eventRepo.Create(ctx, cecReq)
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}
		}
		err = eu.eventRepo.Delete(ctx, relatedEventID)
		if err != nil {
			return fmt.Errorf("failed to delete event: %w", err)
		}
		return nil
	})
	return err
}

func (eu *EventUseCase) deletingCorrectEvent(relatedEventMoney int32, creditEventCreateReqs []*EventCreateRequest) bool {
	moneySum := int32(0)
	for _, cecReq := range creditEventCreateReqs {
		moneySum += cecReq.Money
	}

	return moneySum == relatedEventMoney
}

func (eu *EventUseCase) getEventIDFromSplitBaseTag(ctx context.Context, splitBaseTagName string,
	date time.Time) (int64, error) {
	from := date.AddDate(0, -2, -5)
	events, err := eu.eventRepo.ListOutcomesWithTags(ctx, []string{splitBaseTagName}, from, date)
	if err != nil {
		return 0, err
	}
	return events[len(events)-1].GetID(), nil
}

func (eu *EventUseCase) Split(ctx context.Context, eventID int64, splitBaseTagName string,
	date time.Time, money int32, desc string) error {
	if eventID == -1 {
		var err error
		eventID, err = eu.getEventIDFromSplitBaseTag(ctx, splitBaseTagName, date)
		if err != nil {
			return err
		}
		slog.Info("Auto detected.", "eventID", eventID)
	}

	event, err := eu.eventRepo.GetWithoutTags(ctx, eventID)
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

	err = eu.tx.Do(ctx, func(ctx context.Context) error {
		// Update the existing event.
		if event.GetMoney() == money {
			err = eu.eventRepo.Delete(ctx, eventID)
			if err != nil {
				return fmt.Errorf("failed to delete event: %w", err)
			}
		} else {
			err = eu.eventRepo.UpdateMoney(ctx, eventID, event.GetMoney()-money)
			if err != nil {
				return fmt.Errorf("failed to update event money: %w", err)
			}
		}

		// Insert a new event.
		_, err = eu.eventRepo.Create(ctx, &EventCreateRequest{
			Date:  date,
			Money: money,
			Desc:  model.FormatDesc(desc),
		})
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}
		return nil
	})
	return err
}

func (eu *EventPresentUseCase) PresentOutcomes(ctx context.Context, tagNames []string, from, to time.Time) error {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.ListOutcomes(ctx, from, to)
		if err != nil {
			return fmt.Errorf("failed to list outcome events: %w", err)
		}
	} else {
		events, err = eu.eventRepo.ListOutcomesWithTags(ctx, tagNames, from, to)
		if err != nil {
			return fmt.Errorf("failed to list outcome events with tags: %w", err)
		}
	}

	eu.eventPresenter.Present(events)
	return nil
}

func (eu *EventPresentUseCase) PresentAll(ctx context.Context, tagNames []string, from, to time.Time) error {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.List(ctx, from, to)
		if err != nil {
			return fmt.Errorf("failed to list all events: %w", err)
		}
	} else {
		events, err = eu.eventRepo.ListWithTags(ctx, tagNames, from, to)
		if err != nil {
			return fmt.Errorf("failed to list all events with tags: %w", err)
		}
	}

	eu.eventPresenter.Present(events)
	return nil
}

func (etmu *EventTagMapUsecase) AddTag(ctx context.Context, eventID int64, tagNames []string) error {
	for _, tagName := range tagNames {
		err := etmu.etmRepo.Map(ctx, eventID, tagName)
		if err != nil {
			return fmt.Errorf("failed to remove tag: %w", err)
		}
	}
	return nil
}

func (etmu *EventTagMapUsecase) RemoveTag(ctx context.Context, eventID int64, tagName string) error {
	err := etmu.etmRepo.Unmap(ctx, eventID, tagName)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}
	return err
}

func (apu *ApplyPatternUseCase) ApplyPattern(ctx context.Context, from, to time.Time) error {
	err := apu.tx.Do(ctx, func(ctx context.Context) error {
		events, err := apu.eventRepo.List(ctx, from, to)
		if err != nil {
			return err
		}

		patterns, err := apu.patternRepo.List(ctx)
		if err != nil {
			return err
		}

		for _, event := range events {
			for _, pattern := range patterns {
				if strings.Contains(event.GetDesc(), pattern.GetKey()) {
					for _, tagName := range pattern.GetTagNames() {
						err = apu.etmRepo.Map(ctx, event.GetID(), tagName)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})
	return err
}
