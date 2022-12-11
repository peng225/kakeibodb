package mysql_client

import (
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/db_client"
	"log"
	"reflect"
	"strconv"
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

func (mc *MySQLClient) Insert(table string, withID bool, data []any) error {
	if data == nil || len(data) == 0 {
		return errors.New("empty data.")
	}
	queryString := "insert into " + table + " VALUES(?"
	if withID {
		queryString += ",?"
	}
	queryString += strings.Repeat(",?", len(data)-1) + ")"
	stmtIns, err := mc.db.Prepare(queryString)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	insertData := make([]any, 0)
	if withID {
		insertData = append(insertData, 0)
	}
	insertData = append(insertData, data...)
	_, err = stmtIns.Exec(insertData...)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MySQLClient) SelectPaymentEvent(from, to string) {
	queryStr := fmt.Sprintf("select %s.*, group_concat(%s.name separator ', ') as tags from %s left outer join %s on %s.id = %s.event_id left outer join %s on %s.id = %s.tag_id where (event.dt between '%s' and '%s') and (event.money < 0) group by %s.id order by event.dt;",
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

func (mc *MySQLClient) SelectPaymentEventWithAllTags(tags []string, from, to string) {
	singleQuotedTags := singleQuoteEachString(tags)
	queryStr := fmt.Sprintf("select %s.*, group_concat(%s.name separator ', ') as tags from %s left outer join %s on %s.id = %s.event_id left outer join %s on %s.id = %s.tag_id where (event.dt between '%s' and '%s') and (tag.name in (%s)) and (event.money < 0) group by %s.id order by event.dt;",
		db_client.EventTableName, db_client.TagTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.EventTableName, db_client.MapTableName,
		db_client.TagTableName, db_client.TagTableName, db_client.MapTableName,
		from, to, strings.Join(singleQuotedTags, ","),
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

func (mc *MySQLClient) Select(table string, param any) ([]string, [][]string, error) {
	queryStr := fmt.Sprintf("select * from %s", table)

	if param != nil {
		st := reflect.TypeOf(param)
		sv := reflect.ValueOf(param)
		firstField := true
		for i := 0; i < st.NumField(); i++ {
			ft := st.Field(i)
			fv := sv.Field(i)
			var valueStr string
			var ok bool
			switch fv.Kind() {
			case reflect.String:
				valueStr, ok = fv.Interface().(string)
				if !ok {
					return nil, nil, fmt.Errorf("type conversion failed. fv.Kind = %v", fv.Kind())
				}
				if valueStr == "" {
					continue
				}
				valueStr = "'" + valueStr + "'"
			case reflect.Int:
				valueInt, ok := fv.Interface().(int)
				if !ok {
					return nil, nil, fmt.Errorf("type conversion failed. fv.Kind = %v", fv.Kind())
				}
				if valueInt == 0 {
					continue
				}
				valueStr = strconv.Itoa(valueInt)
			}
			if firstField {
				queryStr += " where "
				firstField = false
			} else {
				queryStr += " and "
			}
			queryStr += fmt.Sprintf("%s = %s", ft.Tag.Get("colName"), valueStr)
		}
	}

	rows, err := mc.db.Query(queryStr)
	if err != nil {
		return nil, nil, err
	}

	header, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	entries := [][]string{}
	for rows.Next() {
		entry := make([]string, len(header))
		switch len(entry) {
		case 2:
			err = rows.Scan(&entry[0], &entry[1])
		case 3:
			err = rows.Scan(&entry[0], &entry[1], &entry[2])
		case 4:
			err = rows.Scan(&entry[0], &entry[1], &entry[2], &entry[3])
		case 5:
			err = rows.Scan(&entry[0], &entry[1], &entry[2], &entry[3], &entry[4])
		default:
			log.Fatalf("Invalid number of columns: %d", len(entry))
		}
		if err != nil {
			return nil, nil, err
		}
		entries = append(entries, entry)
	}
	return header, entries, nil
}

func (mc *MySQLClient) DeleteByID(table string, id int) error {
	stmtIns, err := mc.db.Prepare("delete from " + table + " where id = ?")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(id)
	if err != nil {
		return err
	}
	return nil
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
		if errors.Is(err, sql.ErrNoRows) {
			log.Fatal(err)
		}
		money = 0
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
