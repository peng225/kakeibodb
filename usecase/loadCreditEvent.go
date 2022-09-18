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

		log.Printf("insert value (%v, %v, %v, %v))\n", relatedBankEventID, date, money, desc)
		leh.dbClient.InsertCreditEvent(relatedBankEventID, date, money, desc)
	}
}
