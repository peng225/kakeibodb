package usecase

import (
	"fmt"
	"kakeibodb/db_client"
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
