package usecase

import "kakeibodb/db_client"

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

	th.dbClient.InsertTag(name)
}

func (th *TagHandler) DeleteTag(id int) {
	th.dbClient.Open(db_client.DBName, "shinya")
	defer th.dbClient.Close()

	th.dbClient.DeleteTag(id)
}
