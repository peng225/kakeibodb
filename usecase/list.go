package usecase

import (
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

	lh.dbClient.SelectTagAll()
}
