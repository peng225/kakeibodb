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
	rows, err := mc.db.Query(fmt.Sprintf("select * from %s", db_client.EventTableName))
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
		err := rows.Scan(&id, &date, &money, &desc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\t%v\t%v\t%v\n", id, date, money, desc)
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
