package mysql_client

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	db *sql.DB
}

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

func (mc *MySQLClient) InsertTag(tag string) {
	stmtIns, err := mc.db.Prepare("insert into event VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(0, tag)
	if err != nil {
		log.Fatal(err)
	}
}

func (mc *MySQLClient) InsertMap(eventID, tagID int) {
	stmtIns, err := mc.db.Prepare("insert into event VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eventID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}
