package usecase

import (
	"kakeibodb/db_client"
	"kakeibodb/event"
	"log"
	"strconv"
)

type LoadCreditEventHandler struct {
	dbClient db_client.DBClient
}

type creditEvent struct {
	date        string
	money       int
	description string
}

func NewLoadCreditEventHandler(dc db_client.DBClient) *LoadCreditEventHandler {
	return &LoadCreditEventHandler{
		dbClient: dc,
	}
}

func (leh *LoadCreditEventHandler) LoadCreditEventFromFile(file string, relatedBankEventID int) {
	csv := event.NewCSV()
	csv.Open(file)

	leh.dbClient.Open(db_client.DBName, "shinya")
	defer leh.dbClient.Close()

	log.Printf("load from %s\n", file)

	// Skip header
	_ = csv.Read()
	creditEvents := []creditEvent{}
	for {
		event := csv.Read()
		if event == nil {
			break
		}

		date := event[0]
		if date == "" {
			continue
		}
		desc := event[1]
		money, err := strconv.Atoi(event[2])
		if err != nil {
			log.Fatal(err)
		}
		money *= -1

		creditEvents = append(creditEvents, creditEvent{
			date:        date,
			money:       money,
			description: desc,
		})
	}

	if !leh.deletingCorrectEvent(relatedBankEventID, creditEvents) {
		log.Fatalf("deleting invalid event or event not found. ID = %v", relatedBankEventID)
	}
	for _, ce := range creditEvents {
		shortDesc := string([]rune(ce.description)[0:32])
		dup, err := leh.hasDuplicateEvent(ce.date, ce.money, shortDesc)
		if err != nil {
			log.Fatal(err)
		}
		if dup {
			log.Printf("duplicate event found. date = %v, money = %v, desc = %v", ce.date, ce.money, shortDesc)
			continue
		}
		log.Printf("insert value (%v, %v, %v)\n", ce.date, ce.money, string([]rune(ce.description)[0:32]))
		var insertData []any = []any{ce.date, ce.money, string([]rune(ce.description)[0:32])}
		err = leh.dbClient.Insert(db_client.EventTableName, true, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := leh.dbClient.DeleteByID(db_client.EventTableName, relatedBankEventID)
	if err != nil {
		log.Fatal(err)
	}
}

func (leh *LoadCreditEventHandler) deletingCorrectEvent(id int, creditEvents []creditEvent) bool {
	moneySum := 0
	for _, ce := range creditEvents {
		moneySum += ce.money
	}
	tagEntry := db_client.TagEntry{
		ID: id,
	}
	_, entries, err := leh.dbClient.Select(db_client.EventTableName, tagEntry)
	if err != nil {
		log.Fatal(err)
	}
	if len(entries) == 0 {
		return false
	}
	money, err := strconv.Atoi(entries[0][1])
	return moneySum == money
}

func (leh *LoadCreditEventHandler) hasDuplicateEvent(date string, money int, desc string) (bool, error) {
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
