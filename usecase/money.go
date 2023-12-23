package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type MoneyHandler struct {
	dbClient db_client.DBClient
}

func NewMoneyHandler(dc db_client.DBClient) *MoneyHandler {
	dc.Open()
	return &MoneyHandler{
		dbClient: dc,
	}
}

func (mh *MoneyHandler) Close() {
	mh.dbClient.Close()
}

func (mh *MoneyHandler) GetTotalMoney(tags, from, to string) {
	var money int
	if tags == "" {
		money = mh.dbClient.GetMoneySum(from, to)
	} else if (!strings.Contains(tags, "&") && !strings.Contains(tags, "|")) ||
		strings.Contains(tags, "&") {
		tagTokens := strings.Split(tags, "&")
		money = mh.dbClient.GetMoneySumForAllTags(tagTokens, from, to)
	} else {
		tagTokens := strings.Split(tags, "|")
		money = mh.dbClient.GetMoneySumForAnyTags(tagTokens, from, to)
	}
	fmt.Printf("money: %d\n", money)
}

type MoneyAndTagEntry struct {
	money    int
	tagEntry db_client.TagEntry
}

func (mh *MoneyHandler) Rank(from, to string) {
	totalMoney, mtes, err := mh.rank(from, to, "money")
	if err != nil {
		log.Fatal(err)
	}
	p := message.NewPrinter(language.English)
	p.Printf("total: %d\n", totalMoney)
	for _, mte := range mtes {
		p.Printf("%-8s:\t%8d (%f%%)\n", mte.tagEntry.TagName, mte.money,
			float32(100.0)*float32(mte.money)/float32(totalMoney))
	}
}

func (mh *MoneyHandler) rank(from, to, sortKey string) (int, []*MoneyAndTagEntry, error) {
	var totalMoney int
	totalMoney = mh.dbClient.GetMoneySum(from, to)

	_, tagEntries, err := mh.dbClient.Select(db_client.TagTableName, nil)
	if err != nil {
		return 0, nil, err
	}

	mtes := []*MoneyAndTagEntry{}
	for _, te := range tagEntries {
		tagName := te[db_client.TagColName]
		money := mh.dbClient.GetMoneySumForAllTags([]string{tagName}, from, to)

		id, err := strconv.Atoi(te[db_client.TagColID])
		if err != nil {
			return 0, nil, err
		}
		mtes = append(mtes, &MoneyAndTagEntry{money: money, tagEntry: db_client.TagEntry{
			ID:      id,
			TagName: tagName,
		}})
	}
	money := mh.dbClient.GetMoneySumWithoutTag(from, to)
	mtes = append(mtes, &MoneyAndTagEntry{money: money, tagEntry: db_client.TagEntry{
		ID:      -1,
		TagName: "not tagged",
	}})

	err = validateAnalyzeResult(totalMoney, mtes)
	if err != nil {
		return 0, nil, err
	}

	var mtesSort func(i, j int) bool
	switch sortKey {
	case "money":
		mtesSort = func(i, j int) bool { return mtes[i].money > mtes[j].money }
	case "id":
		mtesSort = func(i, j int) bool { return mtes[i].tagEntry.ID < mtes[j].tagEntry.ID }
	default:
		return 0, nil, fmt.Errorf("invalid sortKey: %s", sortKey)
	}

	sort.Slice(mtes, mtesSort)
	return totalMoney, mtes, nil
}

func validateAnalyzeResult(total int, mtes []*MoneyAndTagEntry) error {
	calculatedTotal := 0
	for _, mte := range mtes {
		calculatedTotal += mte.money
	}
	if calculatedTotal < total {
		return fmt.Errorf("validation error. total: %d, calculatedTotal: %d",
			total, calculatedTotal)
	}
	return nil
}

func (mh *MoneyHandler) TimeSeries(from, to string, interval, window int) {
	layout := "2006-01-02"
	startTime, err := time.Parse(layout, from)
	if err != nil {
		log.Fatal(err)
	}
	endTime, err := time.Parse(layout, to)
	if err != nil {
		log.Fatal(err)
	}

	p := message.NewPrinter(language.English)
	first := true
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.AddDate(0, interval, 0) {
		currentFromInTime := currentTime.AddDate(0, -window, 0)
		currentFrom := fmt.Sprintf("%d-%02d-%02d", currentFromInTime.Year(), currentFromInTime.Month(), currentFromInTime.Day())
		currentTo := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
		_, mtes, err := mh.rank(currentFrom, currentTo, "id")
		if err != nil {
			log.Fatal(err)
		}
		if first {
			p.Print("date ")
			for _, mte := range mtes {
				p.Printf("%q ", mte.tagEntry.TagName)
			}
			p.Println("")
			first = false
		}
		p.Printf("%s ", currentTo)
		for _, mte := range mtes {
			p.Printf("%d ", mte.money)
		}
		p.Println("")
	}
}
