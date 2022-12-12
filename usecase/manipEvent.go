package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"log"
	"strconv"
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
		tagID, err := eh.getTagIDFromName(tagName)
		if err != nil {
			log.Fatal(err)
		}
		var insertData []any = []any{eventID, tagID}
		err = eh.dbClient.Insert(db_client.EventToTagTableName, false, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eh *EventHandler) RemoveTag(eventID int, tagName string) {
	eh.dbClient.Open(db_client.DBName, "shinya")
	defer eh.dbClient.Close()

	tagID, err := eh.getTagIDFromName(tagName)
	if err != nil {
		log.Fatal(err)
	}
	eventToTagEntry := db_client.EventToTagEntry{
		EventID: eventID,
		TagID:   tagID,
	}
	err = eh.dbClient.Delete(db_client.EventToTagTableName, eventToTagEntry)
	if err != nil {
		log.Fatal(err)
	}
}

func (eh *EventHandler) getTagIDFromName(tagName string) (int, error) {
	tagEntry := db_client.TagEntry{
		TagName: tagName,
	}
	_, entries, err := eh.dbClient.Select(db_client.TagTableName, tagEntry)
	if err != nil {
		return 0, err
	}
	if len(entries) != 1 {
		return 0, fmt.Errorf("invalid number of entries. len(entries) = %d", len(entries))
	}
	tagID, err := strconv.Atoi(entries[0][0])
	if err != nil {
		return 0, err
	}
	return tagID, nil
}
