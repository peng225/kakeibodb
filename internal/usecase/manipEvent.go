package usecase

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"kakeibodb/internal/db_client"
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
			continue
		}

		var insertData []any = []any{eventID, tagID}
		_, err = eh.dbClient.Insert(db_client.EventToTagTableName, false, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eh *EventHandler) ApplyPattern(from, to string) {
	_, events, err := eh.dbClient.Select(db_client.EventTableName,
		fmt.Sprintf("where %s between '%s' and '%s'",
			db_client.EventColDate, from, to))
	if err != nil {
		log.Fatal(err)
	}

	_, patterns, err := eh.dbClient.Select(db_client.PatternTableName, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range events {
		for _, pattern := range patterns {
			key := pattern[db_client.PatternColKey]
			desc := event[db_client.EventColDescription]
			if strings.Contains(desc, key) {
				// Get tagID from tagName.
				patternID, err := strconv.Atoi(pattern[db_client.PatternColID])
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
					tagID, err := strconv.Atoi(ptt[db_client.PatternToTagColTID])
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
						tagNames = append(tagNames, tag[db_client.TagColName])
					}
				}
				eventID, err := strconv.Atoi(event[db_client.EventColID])
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
		return 0, fmt.Errorf("tag not found. tagName = %s", tagName)
	}
	tagID, err := strconv.Atoi(entries[0][db_client.TagColID])
	if err != nil {
		return 0, err
	}
	return tagID, nil
}

func (eh *EventHandler) GetEventIDFromSplitBaseTag(splitBaseTagName string,
	date string) (int, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Fatal(err)
	}
	y, m, d := t.AddDate(0, -2, -5).Date()
	from := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	entries, err := eh.dbClient.GetPaymentEventWithAllTags([]string{splitBaseTagName}, from, date)
	if err != nil {
		return 0, err
	}
	sort.Slice(entries, func(i, j int) bool {
		layout := "2006-01-02"
		iTime, err := time.Parse(layout, entries[i]["dt"])
		if err != nil {
			log.Fatal(err)
		}
		jTime, err := time.Parse(layout, entries[j]["dt"])
		if err != nil {
			log.Fatal(err)
		}
		return iTime.After(jTime)
	})
	entryID, err := strconv.Atoi(entries[0]["id"])
	if err != nil {
		return 0, err
	}
	return entryID, nil
}

func (eh *EventHandler) Split(eventID int, date string, money int, desc string) {
	eventQuery := db_client.EventEntry{
		ID: eventID,
	}
	_, events, err := eh.dbClient.Select(db_client.EventTableName, eventQuery)
	if err != nil {
		log.Fatal(err)
	}
	if len(events) != 1 {
		log.Fatalf("event not found: eventID = %d", eventID)
	}

	eventMoneyStr := events[0][db_client.EventColMoney]
	eventMoney, err := strconv.Atoi(eventMoneyStr)
	if err != nil {
		log.Fatal(err)
	}
	if eventMoney > money {
		log.Fatalf("eventMoney(%d) should be smaller than or equal to money(%d)", eventMoney, money)
	}

	// Update the existing event.
	condition := make(map[string]string)
	condition[db_client.EventColID] = strconv.Itoa(eventID)
	updateData := make(map[string]string)
	updateData[db_client.EventColMoney] = strconv.Itoa(eventMoney - money)
	err = eh.dbClient.Update(db_client.EventTableName, condition, updateData)
	if err != nil {
		log.Fatal(err)
	}

	// Insert a new event.
	newEventEntry := []any{date, money, desc}
	_, err = eh.dbClient.Insert(db_client.EventTableName, true, newEventEntry)
	if err != nil {
		log.Fatal(err)
	}
}
