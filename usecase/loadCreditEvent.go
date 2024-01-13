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
	dc.Open()
	return &LoadCreditEventHandler{
		dbClient: dc,
	}
}

func (leh *LoadCreditEventHandler) Close() {
	leh.dbClient.Close()
}

func (leh *LoadCreditEventHandler) LoadCreditEventFromFile(file string, relatedBankEventID int) {
	csv := event.NewCSV()
	csv.Open(file)

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
		if len([]rune(ce.description)) >= db_client.EventDescLength {
			ce.description = string([]rune(ce.description)[0:db_client.EventDescLength])
		}
		dup, err := leh.hasDuplicateEvent(ce.date, ce.money, ce.description)
		if err != nil {
			log.Fatal(err)
		}
		if dup {
			log.Printf("duplicate event found. date = %v, money = %v, desc = %v", ce.date, ce.money, ce.description)
			continue
		}
		log.Printf("insert value (%v, %v, %v)\n", ce.date, ce.money, ce.description)
		var insertData []any = []any{ce.date, ce.money, ce.description}
		_, err = leh.dbClient.Insert(db_client.EventTableName, true, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
	eventEntry := db_client.EventEntry{
		ID: relatedBankEventID,
	}
	err := leh.dbClient.Delete(db_client.EventTableName, eventEntry)
	if err != nil {
		log.Fatal(err)
	}
}

func (leh *LoadCreditEventHandler) deletingCorrectEvent(id int, creditEvents []creditEvent) bool {
	moneySum := 0
	for _, ce := range creditEvents {
		moneySum += ce.money
	}
	eventEntry := db_client.EventEntry{
		ID: id,
	}
	_, entries, err := leh.dbClient.Select(db_client.EventTableName, eventEntry)
	if err != nil {
		log.Fatal(err)
	}
	if len(entries) == 0 {
		return false
	}
	money, err := strconv.Atoi(entries[0][db_client.EventColMoney])
	if err != nil {
		log.Fatal(err)
	}
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
