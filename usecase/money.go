package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"log"
	"sort"
	"strconv"
	"strings"

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

func (mh *MoneyHandler) AnalyzeMoney(from, to string) {
	var totalMoney int
	totalMoney = mh.dbClient.GetMoneySum(from, to)
	p := message.NewPrinter(language.English)
	p.Printf("total: %d\n", totalMoney)

	_, tagEntries, err := mh.dbClient.Select(db_client.TagTableName, nil)
	if err != nil {
		log.Fatal(err)
	}

	mtes := []*MoneyAndTagEntry{}
	for _, te := range tagEntries {
		tagName := te[db_client.TagColName]
		money := mh.dbClient.GetMoneySumForAllTags([]string{tagName}, from, to)

		id, err := strconv.Atoi(te[db_client.TagColID])
		if err != nil {
			log.Fatal(err)
		}
		mtes = append(mtes, &MoneyAndTagEntry{money: money, tagEntry: db_client.TagEntry{
			ID:      id,
			TagName: tagName,
		}})
	}

	sort.Slice(mtes, func(i, j int) bool { return mtes[i].money > mtes[j].money })

	for _, mte := range mtes {
		p.Printf("%-8s:\t%8d (%f%%)\n", mte.tagEntry.TagName, mte.money,
			float32(100.0)*float32(mte.money)/float32(totalMoney))
	}
}
