package db_client

type DBClient interface {
	Open(dbName string, user string)
	Close()
	InsertEvent(date string, money int, description string)
	InsertTag(tag string)
	InsertMap(eventID, tagID int)
	SelectEventAll()
	SelectTagAll()
	// Delete(id int, tableName string)
}
