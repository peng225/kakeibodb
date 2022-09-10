package usecase

import (
	"fmt"
	"kakeibodb/event"
	"kakeibodb/mysql_client"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func LoadEventFromFile(file string) {
	bankCSV := event.NewBankCSV()
	bankCSV.Open(file)

	mysqlClient := mysql_client.NewMySQLClient()
	mysqlClient.Open("kakeibo", "shinya")

	// Skip header
	_ = bankCSV.Read()
	for {
		event := bankCSV.Read()
		if event == nil {
			break
		}
		date := event[0]
		decrease := event[1]
		increase := event[2]
		desc := event[3]

		var money int
		var err error
		if (decrease == "" && increase == "") || (decrease != "" && increase != "") {
			log.Fatalf("bad event record. decrease = %s, increase = %s", decrease, increase)
		} else if decrease != "" {
			money, err = strconv.Atoi(decrease)
			if err != nil {
				log.Fatal(err)
			}
			money *= -1
		} else {
			money, err = strconv.Atoi(increase)
			if err != nil {
				log.Fatal(err)
			}
		}
		mysqlClient.InsertEvent(date, money, desc)
	}
}

func LoadEventFromDir(dir string) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		LoadEventFromFile(file)
	}
}
