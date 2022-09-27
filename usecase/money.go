package usecase

import (
	"fmt"
	"kakeibodb/db_client"
	"sort"
	"strings"
)

type MoneyHandler struct {
	dbClient db_client.DBClient
}

func NewMoneyHandler(dc db_client.DBClient) *MoneyHandler {
	return &MoneyHandler{
		dbClient: dc,
	}
}

func (mh *MoneyHandler) GetTotalMoney(tags, from, to string) {
	mh.dbClient.Open(db_client.DBName, "shinya")
	defer mh.dbClient.Close()

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
	mh.dbClient.Open(db_client.DBName, "shinya")
	defer mh.dbClient.Close()

	var totalMoney int
	totalMoney = mh.dbClient.GetMoneySum(from, to)
	fmt.Printf("total: %d\n", totalMoney)

	_, tagEntries := mh.dbClient.SelectTagAll()

	mtes := []*MoneyAndTagEntry{}
	for _, te := range tagEntries {
		money := mh.dbClient.GetMoneySumForAllTags([]string{te.TagName}, from, to)
		mtes = append(mtes, &MoneyAndTagEntry{money: money, tagEntry: te})
	}

	sort.Slice(mtes, func(i, j int) bool { return mtes[i].money > mtes[j].money })

	for _, mte := range mtes {
		fmt.Printf("%-8s:\t%8d (%f%%)\n", mte.tagEntry.TagName, mte.money,
			float32(100.0)*float32(mte.money)/float32(totalMoney))
	}
}
