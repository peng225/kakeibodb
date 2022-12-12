package usecase

import (
	"kakeibodb/db_client"
	"log"
)

type TagHandler struct {
	dbClient db_client.DBClient
}

func NewTagHandler(dc db_client.DBClient) *TagHandler {
	return &TagHandler{
		dbClient: dc,
	}
}

func (th *TagHandler) CreateTag(name string) {
	th.dbClient.Open(db_client.DBName, "shinya")
	defer th.dbClient.Close()

	var insertData []any = []any{name}
	err := th.dbClient.Insert(db_client.TagTableName, true, insertData)
	if err != nil {
		log.Fatal(err)
	}
}

func (th *TagHandler) DeleteTag(id int) {
	th.dbClient.Open(db_client.DBName, "shinya")
	defer th.dbClient.Close()

	tagEntry := db_client.TagEntry{
		ID: id,
	}
	err := th.dbClient.Delete(db_client.TagTableName, tagEntry)
	if err != nil {
		log.Fatal(err)
	}
}
