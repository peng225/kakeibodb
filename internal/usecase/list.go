package usecase

import (
	"fmt"
	"log"
	"strconv"

	"kakeibodb/internal/db_client"
)

type ListHandler struct {
	dbClient db_client.DBClient
}

func NewListHandler(dc db_client.DBClient) *ListHandler {
	dc.Open()
	return &ListHandler{
		dbClient: dc,
	}
}

func (lh *ListHandler) Close() {
	lh.dbClient.Close()
}

func (lh *ListHandler) ListTag() {
	header, tagEntries, err := lh.dbClient.Select(db_client.TagTableName, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, column := range header {
		fmt.Printf("%s\t", column)
	}
	fmt.Println("")

	for _, te := range tagEntries {
		id, err := strconv.Atoi(te[db_client.TagColID])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%2d\t%v\n", id, te[db_client.TagColName])
	}
}

func (lh *ListHandler) ListAllPattern() {
	lh.dbClient.SelectPatternAll()
}
