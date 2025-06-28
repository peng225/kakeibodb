package usecase

import (
	"fmt"
	"kakeibodb/internal/event"
	"kakeibodb/internal/model"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type EventRepository interface {
	Create(event *model.Event) (int64, error)
	Exist(event *model.Event) (bool, error)
	Get(id int32) (*model.Event, error)
	Delete(id int32) error
}

type EventUseCase struct {
	eventRepo EventRepository
}

func NewEventUseCase(eventRepo EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
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
		e := model.NewEvent(*date, money, desc, nil)
		dup, err := eu.eventRepo.Exist(e)
		if err != nil {
			log.Fatal(err)
		}
		if dup {
			log.Printf("duplicate event found. date = %v, money = %v, desc = %v", date, money, desc)
			continue
		}
		log.Printf("create value (%v, %v, %v)\n", date, money, desc)
		_, err = eu.eventRepo.Create(e)
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

func (eu *EventUseCase) LoadCreditFromFile(file string, relatedEventID int32) {
	csv := event.NewCSV()
	csv.Open(file)

	log.Printf("load from %s\n", file)

	relatedEvent, err := eu.eventRepo.Get(relatedEventID)
	if err != nil {
		log.Fatal(err)
	}

	// Skip header
	_ = csv.Read()
	creditEvents := make([]*model.Event, 0)
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

		creditEvents = append(creditEvents, model.NewEvent(*date, money, desc, nil))
	}

	if !eu.deletingCorrectEvent(relatedEvent, creditEvents) {
		log.Fatalf("deleting invalid event. ID = %v", relatedEventID)
	}
	for _, ce := range creditEvents {
		dup, err := eu.eventRepo.Exist(ce)
		if err != nil {
			log.Fatal(err)
		}
		if dup {
			log.Printf("duplicate event found. date = %v, money = %v, desc = %v",
				ce.GetDate(), ce.GetMoney(), ce.GetDesc())
			continue
		}
		log.Printf("create value (%v, %v, %v)\n", ce.GetDate(), ce.GetMoney(), ce.GetDesc())
		_, err = eu.eventRepo.Create(ce)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = eu.eventRepo.Delete(relatedEventID)
	if err != nil {
		log.Fatal(err)
	}
}

func (eu *EventUseCase) deletingCorrectEvent(relatedEvent *model.Event, creditEvents []*model.Event) bool {
	moneySum := int32(0)
	for _, ce := range creditEvents {
		moneySum += ce.GetMoney()
	}

	return moneySum == relatedEvent.GetMoney()
}
