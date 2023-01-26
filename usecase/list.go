package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"log"
	"strconv"
	"strings"
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

func (lh *ListHandler) ListPaymentEvent(tags, from, to string) {
	if tags == "" {
		lh.dbClient.SelectPaymentEvent(from, to)
	} else {
		lh.dbClient.SelectPaymentEventWithAllTags(strings.Split(tags, "&"), from, to)
	}
}

func (lh *ListHandler) ListAllEvent(from, to string) {
	lh.dbClient.SelectEventAll(from, to)
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
		id, err := strconv.Atoi(te[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%2d\t%v\n", id, te[1])
	}
}

func (lh *ListHandler) ListAllPattern() {
	lh.dbClient.SelectPatternAll()
}
