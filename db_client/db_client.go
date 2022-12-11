package db_client

type EventEntry struct {
	ID    int    `colName:"id"`
	Date  string `colName:"dt"`
	Money int    `colName:"money"`
	Desc  string `colName:"description"`
}

type TagEntry struct {
	ID      int    `colName:"id"`
	TagName string `colName:"name"`
}

type DBClient interface {
	Open(dbName string, user string)
	Close()
	Insert(table string, withID bool, data []any) error
	SelectPaymentEvent(from, to string)
	SelectPaymentEventWithAllTags(tags []string, from, to string)
	SelectEventAll(from, to string)
	Select(table string, param any) ([]string, [][]string, error)
	DeleteByID(table string, id int) error
	DeleteMap(eventID, tagID int)
	GetTagIDFromName(tagName string) int
	GetMoneySum(from, to string) int
	GetMoneySumForAllTags(tags []string, from, to string) int
	GetMoneySumForAnyTags(tags []string, from, to string) int
}
