package usecase

import (
	"fmt"
	"kakeibodb/internal/db_client"
	"kakeibodb/internal/event"
	"kakeibodb/internal/model"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type EventRepository interface {
	Create(event *model.Event) (int64, error)
	Exist(event *model.Event) (bool, error)
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
		layout := "2006/1/2"
		date, err := time.Parse(layout, event[0])
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
		if len([]rune(desc)) >= db_client.EventDescLength {
			desc = string([]rune(desc)[0:db_client.EventDescLength])
		}
		e := model.NewEvent(date, money, desc, nil)
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
