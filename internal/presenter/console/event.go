package console

import (
	"fmt"
	"kakeibodb/internal/model"
	"strings"
)

type EventPresenter struct {
}

func NewEventPresenter() *EventPresenter {
	return &EventPresenter{}
}

func (ep *EventPresenter) Present(events []*model.EventWithID) {
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n", "ID", "date      ", "money   ", "description                     ", "tags")
	for _, e := range events {
		tagsToPrint := "NONE"
		if len(e.GetTagNames()) != 0 {
			tagsToPrint = strings.Join(e.GetTagNames(), ",")
		}
		fmt.Printf("%v\t%v\t%8d\t%-32s\t%s\n",
			e.GetID(), e.GetDate().Format("2006-01-02"), e.GetMoney(), e.GetDesc(), tagsToPrint)
	}
}
