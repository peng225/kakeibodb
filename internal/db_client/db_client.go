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

type EventToTagEntry struct {
	EventID int `colName:"event_id"`
	TagID   int `colName:"tag_id"`
}

type PatternEntry struct {
	ID  int    `colName:"id"`
	Key string `colName:"key_string"`
}

type PatternToTagEntry struct {
	PatternID int `colName:"pattern_id"`
	TagID     int `colName:"tag_id"`
}

type DBClient interface {
	Open()
	Close()
	Select(table string, param any) ([]string, []map[string]string, error)
	GetIncomeSum(from, to string) int
	GetOutcomeSum(from, to string) int
	GetOutcomeSumForAllTags(tags []string, from, to string) int
	GetOutcomeSumForAnyTags(tags []string, from, to string) int
	GetOutcomeSumWithoutTag(from, to string) int
}
