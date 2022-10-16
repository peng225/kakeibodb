package db_client

type TagEntry struct {
	ID      int
	TagName string
}

type DBClient interface {
	Open(dbName string, user string)
	Close()
	Insert(table string, withID bool, data []any) error
	SelectEvent(id int) (string, int, string)
	SelectPaymentEvent(from, to string)
	SelectPaymentEventWithAllTags(tags []string, from, to string)
	SelectEventAll(from, to string)
	SelectTagAll() ([]string, []TagEntry)
	DeleteByID(table string, id int) error
	DeleteMap(eventID, tagID int)
	GetTagIDFromName(tagName string) int
	GetMoneySum(from, to string) int
	GetMoneySumForAllTags(tags []string, from, to string) int
	GetMoneySumForAnyTags(tags []string, from, to string) int
}
