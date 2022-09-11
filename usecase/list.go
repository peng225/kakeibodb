package usecase

import (
	"kakeibodb/db_client"
)

type ListHandler struct {
	dbClient db_client.DBClient
}

func NewListHander(dc db_client.DBClient) *ListHandler {
	return &ListHandler{
		dbClient: dc,
	}
}

func (lh *ListHandler) ListEvent() {
	lh.dbClient.Open(db_client.DBName, "shinya")
	defer lh.dbClient.Close()

	lh.dbClient.SelectEventAll()
}

func (lh *ListHandler) ListTag() {
	lh.dbClient.Open(db_client.DBName, "shinya")
	defer lh.dbClient.Close()

	lh.dbClient.SelectTagAll()
}
