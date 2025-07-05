package mysql_client

import (
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/db_client"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	db     *sql.DB
	dbName string
	dbPort int
	user   string
}

var _ db_client.DBClient = (*MySQLClient)(nil)

func NewMySQLClient(dbName string, dbPort int, user string) *MySQLClient {
	return &MySQLClient{
		dbName: dbName,
		dbPort: dbPort,
		user:   user,
	}
}

func (mc *MySQLClient) Open() {
	var err error
	mc.db, err = sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:%d)/%s",
		mc.user, mc.dbPort, mc.dbName))
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

func (mc *MySQLClient) Insert(table string, withID bool, data []any) (int64, error) {
	if data == nil || len(data) == 0 {
		return 0, errors.New("empty data.")
	}
	queryString := "insert into " + table + " VALUES(?"
	if withID {
		queryString += ",?"
	}
	queryString += strings.Repeat(",?", len(data)-1) + ")"
	stmtIns, err := mc.db.Prepare(queryString)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	insertData := make([]any, 0)
	if withID {
		insertData = append(insertData, 0)
	}
	insertData = append(insertData, data...)
	result, err := stmtIns.Exec(insertData...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		// Not all tables have auto-incremented ID column.
		return 0, nil
	}
	return id, nil
}

func (mc *MySQLClient) Select(table string, param any) ([]string, []map[string]string, error) {
	queryStr := fmt.Sprintf("select * from %s", table)

	if param != nil {
		sv := reflect.ValueOf(param)
		if sv.Kind() == reflect.String {
			queryBody, ok := sv.Interface().(string)
			if !ok {
				return nil, nil, fmt.Errorf("type conversion failed. sv.Kind = %v", sv.Kind())
			}
			queryStr += " " + queryBody
		} else if sv.Kind() == reflect.Struct {
			st := reflect.TypeOf(param)
			firstField := true
			for i := 0; i < sv.NumField(); i++ {
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
				ft := st.Field(i)
				queryStr += fmt.Sprintf("%s = %s", ft.Tag.Get("colName"), valueStr)
			}
		} else {
			return nil, nil, fmt.Errorf("invalid param type. sv.Kind = %v", sv.Kind())
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

	entries := make([]map[string]string, 0)
	for rows.Next() {
		entry := make(map[string]string)
		scanData := make([]string, len(header))
		switch len(scanData) {
		case 2:
			err = rows.Scan(&scanData[0], &scanData[1])
		case 3:
			err = rows.Scan(&scanData[0], &scanData[1], &scanData[2])
		case 4:
			err = rows.Scan(&scanData[0], &scanData[1], &scanData[2], &scanData[3])
		case 5:
			err = rows.Scan(&scanData[0], &scanData[1], &scanData[2], &scanData[3], &scanData[4])
		default:
			log.Fatalf("Invalid number of columns: %d", len(entry))
		}
		if err != nil {
			return nil, nil, err
		}
		for i, d := range scanData {
			entry[header[i]] = d
		}
		entries = append(entries, entry)
	}
	return header, entries, nil
}

func (mc *MySQLClient) GetIncomeSum(from, to string) int {
	row := mc.db.QueryRow(fmt.Sprintf("select sum(%s.money) from %s where (%s.dt between '%s' and '%s') and (%s.money > 0);",
		db_client.EventTableName, db_client.EventTableName,
		db_client.EventTableName, from, to,
		db_client.EventTableName))

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

func (mc *MySQLClient) GetOutcomeSum(from, to string) int {
	row := mc.db.QueryRow(fmt.Sprintf("select -sum(%s.money) from %s where (%s.dt between '%s' and '%s') and (%s.money < 0);",
		db_client.EventTableName, db_client.EventTableName,
		db_client.EventTableName, from, to,
		db_client.EventTableName))

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

func singleQuoteEachString(tags []string) []string {
	resultTags := make([]string, len(tags))
	for i, tag := range tags {
		resultTags[i] = "'" + tag + "'"
	}
	return resultTags
}

func (mc *MySQLClient) GetOutcomeSumForAllTags(tags []string, from, to string) int {
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

func (mc *MySQLClient) GetOutcomeSumForAnyTags(tags []string, from, to string) int {
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

func (mc *MySQLClient) GetOutcomeSumWithoutTag(from, to string) int {
	queryStr := fmt.Sprintf("select -sum(money) from event left outer join event_to_tag on event.id = event_to_tag.event_id where tag_id is NULL and money < 0 and dt between '%s' and '%s';",
		from, to)
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
