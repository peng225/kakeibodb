package usecase

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"kakeibodb/internal/db_client"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type AnalysisHandler struct {
	dbClient db_client.DBClient
}

func NewAnalysisHandler(dc db_client.DBClient) *AnalysisHandler {
	dc.Open()
	return &AnalysisHandler{
		dbClient: dc,
	}
}

func (ah *AnalysisHandler) Close() {
	ah.dbClient.Close()
}

type MoneyAndTagEntry struct {
	money    int
	tagEntry db_client.TagEntry
}

func (ah *AnalysisHandler) rank(from, to, sortKey string) (int, []*MoneyAndTagEntry, error) {
	totalOutcome := ah.dbClient.GetOutcomeSum(from, to)
	if totalOutcome == 0 {
		return 0, nil, nil
	}

	_, tagEntries, err := ah.dbClient.Select(db_client.TagTableName, nil)
	if err != nil {
		return 0, nil, err
	}

	mtes := []*MoneyAndTagEntry{}
	for _, te := range tagEntries {
		tagName := te[db_client.TagColName]
		outcome := ah.dbClient.GetOutcomeSumForAllTags([]string{tagName}, from, to)

		id, err := strconv.Atoi(te[db_client.TagColID])
		if err != nil {
			return 0, nil, err
		}
		mtes = append(mtes, &MoneyAndTagEntry{money: outcome, tagEntry: db_client.TagEntry{
			ID:      id,
			TagName: tagName,
		}})
	}
	outcome := ah.dbClient.GetOutcomeSumWithoutTag(from, to)
	mtes = append(mtes, &MoneyAndTagEntry{money: outcome, tagEntry: db_client.TagEntry{
		ID:      -1,
		TagName: "not tagged",
	}})

	err = validateAnalyzeResult(totalOutcome, mtes)
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
	return totalOutcome, mtes, nil
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

func (ah *AnalysisHandler) TimeSeries(from, to string, interval, window, top int) {
	layout := "2006-01-02"
	startTime, err := time.Parse(layout, from)
	if err != nil {
		log.Fatal(err)
	}
	endTime, err := time.Parse(layout, to)
	if err != nil {
		log.Fatal(err)
	}
	if top < 1 {
		log.Fatal(`"top" must be larger than 0`)
	}

	topIDs, err := ah.getTopIDs(from, to, window, top)
	if err != nil {
		log.Fatal(err)
	}

	p := message.NewPrinter(language.English)
	first := true
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.AddDate(0, interval, 0) {
		currentFromInTime := currentTime.AddDate(0, -window, 0)
		currentFrom := fmt.Sprintf("%d-%02d-%02d", currentFromInTime.Year(), currentFromInTime.Month(), currentFromInTime.Day())
		currentTo := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
		totalOutcome, mtes, err := ah.rank(currentFrom, currentTo, "id")
		if err != nil {
			log.Fatal(err)
		}
		totalIncome := ah.dbClient.GetIncomeSum(currentFrom, currentTo)
		if first {
			p.Print(`"date" "totalIncome" "totalOutcome" `)
			for _, mte := range mtes {
				if _, ok := topIDs[mte.tagEntry.ID]; !ok {
					continue
				}
				p.Printf("%q ", mte.tagEntry.TagName)
			}
			p.Println("")
			first = false
		}
		p.Printf("%s %d %d ", currentTo, totalIncome, totalOutcome)
		for _, mte := range mtes {
			if _, ok := topIDs[mte.tagEntry.ID]; !ok {
				continue
			}
			p.Printf("%d ", mte.money)
		}
		p.Println("")
	}
}

func (ah *AnalysisHandler) getTopIDs(from, to string, window, top int) (map[int]struct{}, error) {
	layout := "2006-01-02"
	fromTime, err := time.Parse(layout, from)
	if err != nil {
		log.Fatal(err)
	}
	fromMinusWindowInTime := fromTime.AddDate(0, -window, 0)
	fromMinusWindow := fmt.Sprintf("%d-%02d-%02d", fromMinusWindowInTime.Year(), fromMinusWindowInTime.Month(), fromMinusWindowInTime.Day())
	_, mtes, err := ah.rank(fromMinusWindow, to, "money")
	if err != nil {
		return nil, err
	}

	topIDs := make(map[int]struct{}, 0)
	for i := 0; i < min(top, len(mtes)); i++ {
		topIDs[mtes[i].tagEntry.ID] = struct{}{}
	}
	return topIDs, nil
}
