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

func (mc *MySQLClient) SelectEvent(id int) (string, int, string) {
	queryStr := fmt.Sprintf("select * from event where id = %d", id)
	row := mc.db.QueryRow(queryStr)

	// Print body
	var tmpID int
	var date string
	var money int
	var desc string
	err := row.Scan(&tmpID, &date, &money, &desc)
	if err != nil {
		log.Fatal(err)
	}
	return date, money, desc
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

func (mc *MySQLClient) DeleteEvent(id int) {
	stmtIns, err := mc.db.Prepare("delete from event where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(id)
	if err != nil {
		log.Fatal(err)
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
