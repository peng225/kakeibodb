package mysql_client

import (
	"database/sql"
	"fmt"
	"kakeibodb/db_client"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	db *sql.DB
}

var _ db_client.DBClient = (*MySQLClient)(nil)

func NewMySQLClient() *MySQLClient {
	return &MySQLClient{}
}

func (mc *MySQLClient) Open(dbName string, user string) {
	var err error
	mc.db, err = sql.Open("mysql", user+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	mc.db.SetConnMaxLifetime(time.Minute * 3)
	mc.db.SetMaxOpenConns(2)
	mc.db.SetMaxIdleConns(2)

	err = mc.db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) Close() {
	mc.db.Close()
}

func (mc *MySQLClient) InsertEvent(date string, money int, description string) {
	stmtIns, err := mc.db.Prepare("insert into event VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(0, date, money, description)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) InsertTag(name string) {
	stmtIns, err := mc.db.Prepare("insert into tag VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(0, name)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) InsertMap(eventID, tagID int) {
	stmtIns, err := mc.db.Prepare("insert into event_to_tag VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eventID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) InsertCreditMap(creditEventID, tagID int) {
	stmtIns, err := mc.db.Prepare("insert into credit_event_to_tag VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(creditEventID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) SelectEventAll(from, to string) {
	queryStr := fmt.Sprintf("select %s.*, group_concat(%s.name separator ', ') as tags from %s left outer join %s on %s.id = %s.event_id left outer join %s on %s.id = %s.tag_id where (event.dt between '%s' and '%s') group by %s.id order by event.dt;",
		db_client.EventTableName, db_client.TagTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.TagTableName, db_client.TagTableName, db_client.MapTableName,
		from, to,
		db_client.EventTableName)
	rows, err := mc.db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Print header
	for _, column := range columns {
		fmt.Printf("%s\t", column)
	}
	fmt.Println("")

	// Print body
	for rows.Next() {
		var id int
		var date string
		var money int
		var desc string
		var tags *string
		err := rows.Scan(&id, &date, &money, &desc, &tags)
		if err != nil {
			log.Fatal(err)
		}
		if tags == nil {
			tmpTags := "NULL"
			tags = &tmpTags
		}
		fmt.Printf("%v\t%v\t%8d\t%-32s\t%s\n", id, date, money, desc, *tags)
	}
}

func (mc *MySQLClient) SelectEventAllWithCredit(from, to string) {
	queryStr := fmt.Sprintf(`select event_with_tag.id, event_with_tag.dt, event_with_tag.money, event_with_tag.description, credit_event_with_tag.id as credit_id, credit_event_with_tag.money as credit_money, credit_event_with_tag.description as credit_description, event_with_tag.tags, credit_event_with_tag.tags as credit_tags from (select event.*, group_concat(tag.name separator ', ') as tags from event left outer join event_to_tag on event.id = event_to_tag.event_id left outer join tag on tag.id = event_to_tag.tag_id where (event.dt between "%s" and "%s") group by event.id) as event_with_tag left outer join (select credit_event.*, group_concat(tag.name separator ', ') as tags from credit_event left outer join credit_event_to_tag on credit_event.id = credit_event_to_tag.credit_event_id left outer join tag on tag.id = credit_event_to_tag.tag_id group by credit_event.id order by credit_event.dt) as credit_event_with_tag on event_with_tag.id = credit_event_with_tag.event_id order by event_with_tag.dt;`,
		from, to)
	rows, err := mc.db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Print header
	for _, column := range columns {
		fmt.Printf("%s\t", column)
	}
	fmt.Println("")

	// Print body
	for rows.Next() {
		var id int
		var date string
		var money int
		var desc string
		var creditID *int
		var creditMoney *int
		var creditDesc *string
		var tags *string
		var creditTags *string
		err := rows.Scan(&id, &date, &money, &desc, &creditID, &creditMoney, &creditDesc, &tags, &creditTags)
		if err != nil {
			log.Fatal(err)
		}
		if tags == nil {
			tmpTags := "NULL"
			tags = &tmpTags
		}
		if creditTags == nil {
			tmpCreditTags := "NULL"
			creditTags = &tmpCreditTags
		}
		if creditID == nil {
			fmt.Printf("%v\t%v\t%8d\t%-32s\tNULL\tNULL\tNULL\t%s\t%s\n", id, date, money, desc, *tags, *creditTags)
		} else {
			fmt.Printf("%v\t%v\t%8d\t%-32s\t%v\t%v\t%-32s\t%s\t%s\n", id, date, money, desc, *creditID, *creditMoney, *creditDesc, *tags, *creditTags)
		}
	}
}

func (mc *MySQLClient) SelectTagAll() {
	rows, err := mc.db.Query(fmt.Sprintf("select * from %s", db_client.TagTableName))
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Print header
	for _, column := range columns {
		fmt.Printf("%s\t", column)
	}
	fmt.Println("")

	// Print body
	for rows.Next() {
		var id int
		var tagName string
		err := rows.Scan(&id, &tagName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\t%v\n", id, tagName)
	}
}

func (mc *MySQLClient) DeleteTag(id int) {
	stmtIns, err := mc.db.Prepare("delete from tag where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) DeleteMap(eventID, tagID int) {
	stmtIns, err := mc.db.Prepare("delete from event_to_tag where event_id = ? and tag_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eventID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) DeleteCreditMap(creditEventID, tagID int) {
	stmtIns, err := mc.db.Prepare("delete from credit_event_to_tag where credit_event_id = ? and tag_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(creditEventID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) GetTagIDFromName(tagName string) int {
	row := mc.db.QueryRow(fmt.Sprintf("select %s.id from %s where %s.name = '%s'",
		db_client.TagTableName, db_client.TagTableName, db_client.TagTableName, tagName))

	var tagID int
	err := row.Scan(&tagID)
	if err != nil {
		log.Fatal(err)
	}
	return tagID
}

func (mc *MySQLClient) GetMoneySum(from, to string) int {
	row := mc.db.QueryRow(fmt.Sprintf("select -sum(%s.money) from %s where (%s.dt between '%s' and '%s') and (%s.money < 0);",
		db_client.EventTableName, db_client.EventTableName,
		db_client.EventTableName, from, to,
		db_client.EventTableName))

	var money int
	err := row.Scan(&money)
	if err != nil {
		log.Fatal(err)
	}
	return money
}

func singleQuoteEachString(tags []string) []string {
	resultTags := make([]string, len(tags))
	for i, tag := range tags {
		resultTags[i] = "'" + tag + "'"
	}
	return resultTags
}

func (mc *MySQLClient) GetMoneySumForAllTags(tags []string, from, to string) int {
	singleQuotedTags := singleQuoteEachString(tags)
	queryStr := fmt.Sprintf("select sum(matched_money.tmp_money) from (select -max(money) as tmp_money from ((event inner join event_to_tag on event.id = event_to_tag.event_id) inner join tag on tag.id = event_to_tag.tag_id) where tag.name in (%s) and event.money < 0 and (event.dt between '%s' and '%s') group by event.id having count(event.id) = %d) as matched_money",
		strings.Join(singleQuotedTags, ","), from, to, len(tags))
	row := mc.db.QueryRow(queryStr)

	var money int
	err := row.Scan(&money)
	if err != nil {
		log.Fatal(err)
	}
	return money
}

func (mc *MySQLClient) GetMoneySumForAnyTags(tags []string, from, to string) int {
	singleQuotedTags := singleQuoteEachString(tags)
	queryStr := fmt.Sprintf("select sum(matched_money.tmp_money) from (select -max(money) as tmp_money from ((event inner join event_to_tag on event.id = event_to_tag.event_id) inner join tag on tag.id = event_to_tag.tag_id) where tag.name in (%s) and event.money < 0 and (event.dt between '%s' and '%s') group by event.id) as matched_money",
		strings.Join(singleQuotedTags, ","), from, to)
	row := mc.db.QueryRow(queryStr)

	var money int
	err := row.Scan(&money)
	if err != nil {
		log.Fatal(err)
	}
	return money
}

func (mc *MySQLClient) InsertCreditEvent(relatedBankEventID int, date string, money int, description string) {
	stmtIns, err := mc.db.Prepare("insert into credit_event VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(0, relatedBankEventID, date, money, description)
	if err != nil {
		log.Fatal(err)
	}
}
