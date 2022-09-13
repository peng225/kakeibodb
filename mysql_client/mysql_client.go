package mysql_client

import (
	"database/sql"
	"fmt"
	"kakeibodb/db_client"
	"log"
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

func (mc *MySQLClient) SelectEventAll() {
	queryStr := fmt.Sprintf("select %s.*, group_concat(%s.name separator ', ') as tags from %s left outer join %s on %s.id = %s.event_id left outer join %s on %s.id = %s.tag_id group by %s.id order by event.dt;",
		db_client.EventTableName, db_client.TagTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.TagTableName, db_client.TagTableName, db_client.MapTableName,
		db_client.EventTableName)
	fmt.Println(queryStr)
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
		var tags any
		err := rows.Scan(&id, &date, &money, &desc, &tags)
		if err != nil {
			log.Fatal(err)
		}
		if tags == nil {
			fmt.Printf("%v\t%v\t%8d\t%-32s\tNULL\n", id, date, money, desc)
		} else {
			fmt.Printf("%v\t%v\t%8d\t%-32s\t%s\n", id, date, money, desc, tags)
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
