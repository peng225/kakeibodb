package usecase

import (
	"slices"
	"time"
)

type AnalysisUseCase struct {
	eventRepo         EventRepository
	analysisPresenter AnalysisPresenter
}

type AnalysisPresenter interface {
	Present(report *TimeSeriesReport)
}

func NewAnalysisUseCase(eventRepo EventRepository, analysisPresenter AnalysisPresenter) *AnalysisUseCase {
	return &AnalysisUseCase{
		eventRepo:         eventRepo,
		analysisPresenter: analysisPresenter,
	}
}

type tagNameAndMoney struct {
	tagName string
	money   int32
}

func (au *AnalysisUseCase) GetMoneySumGroupedByTagName(from, to time.Time) (map[string]int32, error) {
	events, err := au.eventRepo.ListOutcomes(from, to)
	if err != nil {
		return nil, err
	}

	moneySumByTagName := make(map[string]int32)
	for _, e := range events {
		if len(e.GetTagNames()) == 0 {
			moneySumByTagName["NONE"] += e.GetMoney()
			continue
		}
		for _, tagName := range e.GetTagNames() {
			moneySumByTagName[tagName] += e.GetMoney()
		}
	}
	return moneySumByTagName, nil
}

func GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow [](map[string]int32), top int) []string {
	highlyRankedTagNameAndMoneyList := make([]*tagNameAndMoney, 0)
	for _, msgbtForEveryWindow := range msGroupedByTagNameForEveryWindow {
		for tagName, moneySum := range msgbtForEveryWindow {
			if len(highlyRankedTagNameAndMoneyList) < top {
				highlyRankedTagNameAndMoneyList = append(highlyRankedTagNameAndMoneyList,
					&tagNameAndMoney{tagName, moneySum})
			} else
			// Be careful that money and moneySum are negative.
			if highlyRankedTagNameAndMoneyList[len(highlyRankedTagNameAndMoneyList)-1].money < moneySum {
				continue
			} else {
				i := slices.IndexFunc(highlyRankedTagNameAndMoneyList, func(a *tagNameAndMoney) bool {
					return a.tagName == tagName
				})
				if i != -1 {
					// Be careful that money and moneySum are negative.
					if moneySum < highlyRankedTagNameAndMoneyList[i].money {
						highlyRankedTagNameAndMoneyList[i].money = moneySum
					}
				} else {
					insertIndex := -1
					for i, tam := range highlyRankedTagNameAndMoneyList {
						// Be careful that money and moneySum are negative.
						if moneySum < tam.money {
							insertIndex = i
							break
						}
					}
					if insertIndex == -1 {
						panic("BUG: insertIndex should not be -1.")
					}
					// insertIndex is not used here.
					// Ideally, we should insert a new entry at insetIndex.
					// However, inserting an item in slice is slow.
					// Instead, the last item, which has the lowest score,
					// is replaced with a new one and the slice is sorted later.
					highlyRankedTagNameAndMoneyList[len(highlyRankedTagNameAndMoneyList)-1] = &tagNameAndMoney{tagName, moneySum}
				}
			}
			slices.SortFunc(highlyRankedTagNameAndMoneyList, func(a, b *tagNameAndMoney) int {
				if a.money < b.money {
					return -1
				} else if a.money > b.money {
					return 1
				}
				return 0
			})
		}
	}
	highlyRankedTagNames := make([]string, len(highlyRankedTagNameAndMoneyList))
	for i, tam := range highlyRankedTagNameAndMoneyList {
		highlyRankedTagNames[i] = tam.tagName
	}
	return highlyRankedTagNames
}

func (au *AnalysisUseCase) getMoneyTotal(from, to time.Time) (int32, int32, error) {
	events, err := au.eventRepo.List(from, to)
	if err != nil {
		return 0, 0, err
	}
	totalIncome := int32(0)
	totalOutcome := int32(0)
	for _, e := range events {
		if e.GetMoney() >= 0 {
			totalIncome += e.GetMoney()
		} else {
			totalOutcome += e.GetMoney()
		}
	}
	return totalIncome, totalOutcome, nil
}

type TimeSeriesReportItem struct {
	Date                     time.Time
	TotalOutcome             int32
	TotalIncome              int32
	HighlyRankedOutcomeItems []int32
}

type TimeSeriesReport struct {
	HeaderItems []string
	Items       []*TimeSeriesReportItem
}

func (au *AnalysisUseCase) TimeSeries(from, to time.Time, interval, window, top int) error {
	msGroupedByTagNameForEveryWindow := make([](map[string]int32), 0)
	totalIncomeForEveryWindow := make([]int32, 0)
	totalOutcomeForEveryWindow := make([]int32, 0)
	for windowTo := from; windowTo.Before(to.AddDate(0, 0, 1)); windowTo = windowTo.AddDate(0, interval, 0) {
		windowFrom := windowTo.AddDate(0, -window, 0)
		msGroupedByTagName, err := au.GetMoneySumGroupedByTagName(windowFrom, windowTo)
		if err != nil {
			return err
		}
		msGroupedByTagNameForEveryWindow = append(msGroupedByTagNameForEveryWindow, msGroupedByTagName)

		totalIncome, totalOutcome, err := au.getMoneyTotal(windowFrom, windowTo)
		totalIncomeForEveryWindow = append(totalIncomeForEveryWindow, totalIncome)
		totalOutcomeForEveryWindow = append(totalOutcomeForEveryWindow, totalOutcome)
	}

	var report TimeSeriesReport
	report.HeaderItems = []string{"date", "totalIncome", "totalOutcome"}
	report.Items = make([]*TimeSeriesReportItem, len(msGroupedByTagNameForEveryWindow))
	tagNames := GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, top)
	report.HeaderItems = append(report.HeaderItems, tagNames...)
	tmpDate := from
	for i := range len(report.Items) {
		report.Items[i] = &TimeSeriesReportItem{
			Date:         tmpDate,
			TotalIncome:  totalIncomeForEveryWindow[i],
			TotalOutcome: -totalOutcomeForEveryWindow[i],
		}
		tmpDate = tmpDate.AddDate(0, interval, 0)
		report.Items[i].HighlyRankedOutcomeItems = make([]int32, len(tagNames))
		for j, tagName := range tagNames {
			report.Items[i].HighlyRankedOutcomeItems[j] = -msGroupedByTagNameForEveryWindow[i][tagName]
		}
	}
	au.analysisPresenter.Present(&report)
	return nil
}
