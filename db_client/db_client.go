package db_client

type DBClient interface {
	Open(dbName string, user string)
	Close()
	InsertEvent(date string, money int, description string)
	InsertTag(name string)
	InsertMap(eventID, tagID int)
	SelectEvent(id int) (string, int, string)
	SelectEventAll(from, to string)
	SelectTagAll()
	DeleteEvent(id int)
	DeleteTag(id int)
	DeleteMap(eventID, tagID int)
	GetTagIDFromName(tagName string) int
	GetMoneySum(from, to string) int
	GetMoneySumForAllTags(tags []string, from, to string) int
	GetMoneySumForAnyTags(tags []string, from, to string) int
}
