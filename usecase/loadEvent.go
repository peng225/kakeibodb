package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"kakeibodb/event"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type LoadEventHandler struct {
	dbClient db_client.DBClient
}

func NewLoadEventHandler(dc db_client.DBClient) *LoadEventHandler {
	dc.Open()
	return &LoadEventHandler{
		dbClient: dc,
	}
}

func (leh *LoadEventHandler) Close() {
	leh.dbClient.Close()
}

func (leh *LoadEventHandler) LoadEventFromFile(file string) {
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
		date := event[0]
		decrease := event[1]
		increase := event[2]
		desc := event[3]

		var money int
		var err error
		if (decrease == "" && increase == "") || (decrease != "" && increase != "") {
			log.Fatalf("bad event record. decrease = %s, increase = %s", decrease, increase)
		} else if decrease != "" {
			money, err = strconv.Atoi(decrease)
			if err != nil {
				log.Fatal(err)
			}
			money *= -1
		} else {
			money, err = strconv.Atoi(increase)
			if err != nil {
				log.Fatal(err)
			}
		}
		if len([]rune(desc)) >= 32 {
			desc = string([]rune(desc)[0:32])
		}
		dup, err := leh.hasDuplicateEvent(date, money, desc)
		if err != nil {
			log.Fatal(err)
		}
		if dup {
			log.Printf("duplicate event found. date = %v, money = %v, desc = %v", date, money, desc)
			continue
		}
		log.Printf("insert value (%v, %v, %v)\n", date, money, desc)
		var insertData []any = []any{date, money, desc}
		err = leh.dbClient.Insert(db_client.EventTableName, true, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (leh *LoadEventHandler) LoadEventFromDir(dir string) {
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
		leh.LoadEventFromFile(file)
	}
}

func (leh *LoadEventHandler) hasDuplicateEvent(date string, money int, desc string) (bool, error) {
	eventEntry := db_client.EventEntry{
		Date:  date,
		Money: money,
		Desc:  desc,
	}
	_, data, err := leh.dbClient.Select(db_client.EventTableName, eventEntry)
	if err != nil {
		return false, err
	}
	return len(data) != 0, nil
}
