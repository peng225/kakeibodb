package usecase

import (
	"kakeibodb/db_client"
	"log"
)

type TagHandler struct {
	dbClient db_client.DBClient
}

func NewTagHandler(dc db_client.DBClient) *TagHandler {
	dc.Open()
	return &TagHandler{
		dbClient: dc,
	}
}

func (th *TagHandler) Close() {
	th.dbClient.Close()
}

func (th *TagHandler) CreateTag(name string) {
	var insertData []any = []any{name}
	err := th.dbClient.Insert(db_client.TagTableName, true, insertData)
	if err != nil {
		log.Fatal(err)
	}
}

func (th *TagHandler) DeleteTag(id int) {
	tagEntry := db_client.TagEntry{
		ID: id,
	}
	err := th.dbClient.Delete(db_client.TagTableName, tagEntry)
	if err != nil {
		log.Fatal(err)
	}
}
