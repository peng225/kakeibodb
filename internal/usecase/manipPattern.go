package usecase

import (
	"log"

	"kakeibodb/internal/db_client"
)

type PatternHandler struct {
	dbClient db_client.DBClient
}

func NewPatternHandler(dc db_client.DBClient) *PatternHandler {
	dc.Open()
	return &PatternHandler{
		dbClient: dc,
	}
}

func (th *PatternHandler) Close() {
	th.dbClient.Close()
}

func (ph *PatternHandler) AddTag(patternID int, tagNames []string) {
	for _, tagName := range tagNames {
		tagID, err := getTagIDFromName(ph.dbClient, tagName)
		if err != nil {
			log.Fatal(err)
		}

		// Check whether (eventID, tagID) already exists.
		pttEntry := db_client.PatternToTagEntry{
			PatternID: patternID,
			TagID:     tagID,
		}
		_, ptts, err := ph.dbClient.Select(db_client.PatternToTagTableName, pttEntry)
		if len(ptts) != 0 {
			continue
		}

		var insertData []any = []any{patternID, tagID}
		_, err = ph.dbClient.Insert(db_client.PatternToTagTableName, false, insertData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (ph *PatternHandler) RemoveTag(patternID int, tagName string) {
	tagID, err := getTagIDFromName(ph.dbClient, tagName)
	if err != nil {
		log.Fatal(err)
	}
	patternToTagEntry := db_client.PatternToTagEntry{
		PatternID: patternID,
		TagID:     tagID,
	}
	err = ph.dbClient.Delete(db_client.PatternToTagTableName, patternToTagEntry)
	if err != nil {
		log.Fatal(err)
	}
}
