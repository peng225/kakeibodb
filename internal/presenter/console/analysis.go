package console

import (
	"fmt"
	"kakeibodb/internal/usecase"
)

type AnalysisPresenter struct {
}

func NewAnalysisPresenter() *AnalysisPresenter {
	return &AnalysisPresenter{}
}

func printHeader(headerItems []string) {
	for i, hi := range headerItems {
		fmt.Print(hi)
		if i == len(headerItems)-1 {
			fmt.Println("")
		} else {
			fmt.Print("\t")
		}
	}
}

func (ap *AnalysisPresenter) Present(report *usecase.TimeSeriesReport) {
	printHeader(report.HeaderItems)

	for _, item := range report.Items {
		fmt.Printf("%s\t%d\t%d\t",
			item.Date.Format("2006/01/02"), item.TotalIncome, item.TotalOutcome,
		)
		for i, rakedItem := range item.HighlyRankedOutcomeItems {
			fmt.Printf("%d", rakedItem)
			if i == len(item.HighlyRankedOutcomeItems)-1 {
				fmt.Println("")
			} else {
				fmt.Print("\t")
			}
		}
	}
}
