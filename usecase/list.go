package usecase

import (
	"fmt"
	"kakeibodb/db_client"
)

type ListHandler struct {
	dbClient db_client.DBClient
}

func NewListHandler(dc db_client.DBClient) *ListHandler {
	return &ListHandler{
		dbClient: dc,
	}
}

func (lh *ListHandler) ListEvent(from, to string) {
	lh.dbClient.Open(db_client.DBName, "shinya")
	defer lh.dbClient.Close()

	lh.dbClient.SelectEventAll(from, to)
}

func (lh *ListHandler) ListTag() {
	lh.dbClient.Open(db_client.DBName, "shinya")
	defer lh.dbClient.Close()

	header, tagEntries := lh.dbClient.SelectTagAll()
	for _, column := range header {
		fmt.Printf("%s\t", column)
	}
	fmt.Println("")

	for _, te := range tagEntries {
		fmt.Printf("%2d\t%v\n", te.ID, te.TagName)
	}
}
