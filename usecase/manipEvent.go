package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"log"
	"strconv"
	"strings"
)

type EventHandler struct {
	dbClient db_client.DBClient
}

func NewEventHandler(dc db_client.DBClient) *EventHandler {
	dc.Open()
	return &EventHandler{
		dbClient: dc,
	}
}

func (eh *EventHandler) Close() {
	eh.dbClient.Close()
}

func (eh *EventHandler) AddTag(eventID int, tagNames []string) {
	for _, tagName := range tagNames {
		tagID, err := getTagIDFromName(eh.dbClient, tagName)
		if err != nil {
			log.Fatal(err)
		}

		// Check whether (eventID, tagID) already exists.
		ettEntry := db_client.EventToTagEntry{
			EventID: eventID,
			TagID:   tagID,
		}
		_, etts, err := eh.dbClient.Select(db_client.EventToTagTableName, ettEntry)
		if len(etts) != 0 {
			return
		}

		var insertData []any = []any{eventID, tagID}
		err = eh.dbClient.Insert(db_client.EventToTagTableName, false, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eh *EventHandler) RemoveTag(eventID int, tagName string) {
	tagID, err := getTagIDFromName(eh.dbClient, tagName)
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

func (eh *EventHandler) ApplyPattern(from, to string) {
	_, events, err := eh.dbClient.Select(db_client.EventTableName, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, patterns, err := eh.dbClient.Select(db_client.PatternTableName, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range events {
		for _, pattern := range patterns {
			key := pattern[1]
			desc := event[3]
			if strings.Contains(desc, key) {
				// Get tagID from tagName.
				patternID, err := strconv.Atoi(pattern[0])
				if err != nil {
					log.Fatal(err)
				}
				pttEntry := db_client.PatternToTagEntry{
					PatternID: patternID,
				}
				_, ptts, err := eh.dbClient.Select(db_client.PatternToTagTableName, pttEntry)
				if err != nil {
					log.Fatal(err)
				}
				tagNames := make([]string, 0)
				for _, ptt := range ptts {
					// Get tagName from tagID.
					tagID, err := strconv.Atoi(ptt[1])
					if err != nil {
						log.Fatal(err)
					}
					tagEntry := db_client.TagEntry{
						ID: tagID,
					}
					_, tags, err := eh.dbClient.Select(db_client.TagTableName, tagEntry)
					if err != nil {
						log.Fatal(err)
					}

					for _, tag := range tags {
						tagNames = append(tagNames, tag[1])
					}
				}
				eventID, err := strconv.Atoi(event[0])
				if err != nil {
					log.Fatal(err)
				}
				eh.AddTag(eventID, tagNames)
			}
		}
	}
}

func getTagIDFromName(dbClient db_client.DBClient, tagName string) (int, error) {
	tagEntry := db_client.TagEntry{
		TagName: tagName,
	}
	_, entries, err := dbClient.Select(db_client.TagTableName, tagEntry)
	if err != nil {
		return 0, err
	}
	if len(entries) != 1 {
		return 0, fmt.Errorf("invalid number of entries. len(entries) = %d, tagName = %s", len(entries), tagName)
	}
	tagID, err := strconv.Atoi(entries[0][0])
	if err != nil {
		return 0, err
	}
	return tagID, nil
}
