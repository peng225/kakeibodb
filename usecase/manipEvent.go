package usecase

import "kakeibodb/db_client"

type EventHandler struct {
	dbClient db_client.DBClient
}

func NewEventHander(dc db_client.DBClient) *EventHandler {
	return &EventHandler{
		dbClient: dc,
	}
}

func (eh *EventHandler) AddTag(eventID, tagID int) {
	eh.dbClient.Open(db_client.DBName, "shinya")
	defer eh.dbClient.Close()

	eh.dbClient.InsertMap(eventID, tagID)
}

func (eh *EventHandler) RemoveTag(eventID, tagID int) {
	eh.dbClient.Open(db_client.DBName, "shinya")
	defer eh.dbClient.Close()

	eh.dbClient.DeleteMap(eventID, tagID)
}
