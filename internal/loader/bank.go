package loader

import (
	"fmt"
	"io"
	"iter"
	"kakeibodb/internal/model"
	"kakeibodb/internal/usecase"
	"log/slog"
	"strconv"
)

type BankLoader struct {
}

func NewBankLoader() *BankLoader {
	return &BankLoader{}
}

func (bl *BankLoader) Load(r io.Reader) iter.Seq[*usecase.EventCreateRequest] {
	csv := newCSV(r)
	// Skip header
	_ = csv.read()
	return func(yield func(*usecase.EventCreateRequest) bool) {
		for {
			event := csv.read()
			if event == nil {
				return
			}
			date, err := model.ParseDate(event[0])
			if err != nil {
				slog.Error("Failed to parse date.", "err", err.Error())
				return
			}
			decrease := event[1]
			increase := event[2]
			desc := event[3]

			var money int32
			if (decrease == "" && increase == "") || (decrease != "" && increase != "") {
				slog.Error("Bad event record.", "decrease", decrease, "increase", increase)
				return
			} else if decrease != "" {
				tmpMoney, err := strconv.ParseInt(decrease, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf(`Failed to parse "%s" as int`, decrease), "err", err)
					return
				}
				money = -1 * int32(tmpMoney)
			} else {
				tmpMoney, err := strconv.ParseInt(increase, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf(`Failed to parse "%s" as int`, increase), "err", err)
					return
				}
				money = int32(tmpMoney)
			}
			if !yield(&usecase.EventCreateRequest{
				Date:  *date,
				Money: money,
				Desc:  model.FormatDesc(desc),
			}) {
				return
			}
		}
	}
}
