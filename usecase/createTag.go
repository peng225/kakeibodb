package usecase

import "kakeibodb/db_client"

type CreateTagHandler struct {
	dbClient db_client.DBClient
}

func NewCreateTagHander(dc db_client.DBClient) *CreateTagHandler {
	return &CreateTagHandler{
		dbClient: dc,
	}
}

func (cth *CreateTagHandler) CreateTag(name string) {
	cth.dbClient.Open(db_client.DBName, "shinya")
	defer cth.dbClient.Close()

	cth.dbClient.InsertTag(name)
}
