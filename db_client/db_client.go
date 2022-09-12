package db_client

type DBClient interface {
	Open(dbName string, user string)
	Close()
	InsertEvent(date string, money int, description string)
	InsertTag(name string)
	InsertMap(eventID, tagID int)
	SelectEventAll()
	SelectTagAll()
	DeleteTag(id int)
}
