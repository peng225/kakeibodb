package usecase

import (
	"kakeibodb/db_client"
	"log"
)

type EventHandler struct {
	dbClient db_client.DBClient
}

func NewEventHandler(dc db_client.DBClient) *EventHandler {
	return &EventHandler{
		dbClient: dc,
	}
}

func (eh *EventHandler) AddTag(eventID int, tagNames []string) {
	eh.dbClient.Open(db_client.DBName, "shinya")
	defer eh.dbClient.Close()

	for _, tagName := range tagNames {
		tagID := eh.dbClient.GetTagIDFromName(tagName)
		var insertData []any = []any{eventID, tagID}
		err := eh.dbClient.Insert(db_client.MapTableName, false, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eh *EventHandler) RemoveTag(eventID int, tagName string) {
	eh.dbClient.Open(db_client.DBName, "shinya")
	defer eh.dbClient.Close()

	tagID := eh.dbClient.GetTagIDFromName(tagName)
	eh.dbClient.DeleteMap(eventID, tagID)
}
