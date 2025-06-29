package usecase

import (
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

func (lh *ListHandler) ListAllPattern() {
	lh.dbClient.SelectPatternAll()
}
