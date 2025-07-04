package usecase

import (
	"fmt"
	"kakeibodb/internal/event"
	"kakeibodb/internal/model"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	Delete(id int64) error
	ListOutcomes(from, to *time.Time) ([]*model.Event, error)
	ListOutcomesWithTags(tagNames []string, from, to *time.Time) ([]*model.Event, error)
	List(from, to *time.Time) ([]*model.Event, error)
	ListWithTags(tagNames []string, from, to *time.Time) ([]*model.Event, error)
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

func (eu *EventUseCase) LoadFromFile(file string) {
	// FIXME: Don't want to depend on a specific file format.
	csv := event.NewCSV()
	csv.Open(file)

	log.Printf("load from %s\n", file)

	// Skip header
	_ = csv.Read()
	for {
		event := csv.Read()
		if event == nil {
			break
		}
		date, err := model.ParseDate(event[0])
		if err != nil {
			log.Fatal(err)
		}
		decrease := event[1]
		increase := event[2]
		desc := event[3]

		var money int32
		if (decrease == "" && increase == "") || (decrease != "" && increase != "") {
			log.Fatalf("bad event record. decrease = %s, increase = %s", decrease, increase)
		} else if decrease != "" {
			tmpMoney, err := strconv.ParseInt(decrease, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			money = -1 * int32(tmpMoney)
		} else {
			tmpMoney, err := strconv.ParseInt(increase, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			money = int32(tmpMoney)
		}
		log.Printf("create value (%v, %v, %v)\n", date, money, desc)
		_, err = eu.eventRepo.Create(&EventCreateRequest{
			Date:  *date,
			Money: money,
			Desc:  desc,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eu *EventUseCase) LoadFromDir(dir string) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
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
		log.Fatal(err)
	}

	for _, file := range files {
		eu.LoadFromFile(file)
	}
}

func (eu *EventUseCase) LoadCreditFromFile(file string, relatedEventID int64) {
	csv := event.NewCSV()
	csv.Open(file)

	log.Printf("load from %s\n", file)

	relatedEvent, err := eu.eventRepo.GetWithoutTags(relatedEventID)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		desc := event[1]
		tmpMoney, err := strconv.ParseInt(event[2], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		tmpMoney *= -1
		money := int32(tmpMoney)

		creditEventCreateReqs = append(creditEventCreateReqs, &EventCreateRequest{
			Date:  *date,
			Money: money,
			Desc:  desc,
		})
	}

	if !eu.deletingCorrectEvent(relatedEvent.GetMoney(), creditEventCreateReqs) {
		log.Fatalf("sum of credit events does not match original event (ID=%d)", relatedEventID)
	}
	for _, cecReq := range creditEventCreateReqs {
		log.Printf("create value (%v, %v, %v)\n", cecReq.Date, cecReq.Money, cecReq.Desc)
		_, err = eu.eventRepo.Create(cecReq)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = eu.eventRepo.Delete(relatedEventID)
	if err != nil {
		log.Fatal(err)
	}
}

func (eu *EventUseCase) deletingCorrectEvent(relatedEventMoney int32, creditEventCreateReqs []*EventCreateRequest) bool {
	moneySum := int32(0)
	for _, cecReq := range creditEventCreateReqs {
		moneySum += cecReq.Money
	}

	return moneySum == relatedEventMoney
}

func (eu *EventPresentUseCase) PresentOutcomes(tagNames []string, from, to *time.Time) {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.ListOutcomes(from, to)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		events, err = eu.eventRepo.ListOutcomesWithTags(tagNames, from, to)
		if err != nil {
			log.Fatal(err)
		}
	}

	eu.eventPresenter.Present(events)
}

func (eu *EventPresentUseCase) PresentAll(tagNames []string, from, to *time.Time) {
	var events []*model.Event
	var err error
	if len(tagNames) == 0 {
		events, err = eu.eventRepo.List(from, to)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		events, err = eu.eventRepo.ListWithTags(tagNames, from, to)
		if err != nil {
			log.Fatal(err)
		}
	}

	eu.eventPresenter.Present(events)
}

func (etmu *EventTagMapUsecase) AddTag(eventID int64, tagNames []string) {
	for _, tagName := range tagNames {
		err := etmu.etmRepo.Map(eventID, tagName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (etmu *EventTagMapUsecase) RemoveTag(eventID int64, tagName string) {
	err := etmu.etmRepo.Unmap(eventID, tagName)
	if err != nil {
		log.Fatal(err)
	}
}
