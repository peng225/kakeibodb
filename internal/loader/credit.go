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

type CreditLoader struct {
}

func NewCreditLoader() *CreditLoader {
	return &CreditLoader{}
}

func (cl *CreditLoader) Load(r io.Reader) iter.Seq[*usecase.EventCreateRequest] {
	csv := newCSV(r)
	// Skip header
	_ = csv.read()
	return func(yield func(*usecase.EventCreateRequest) bool) {
		for {
			event := csv.read()
			if event == nil {
				return
			}

			if event[0] == "" {
				continue
			}
			date, err := model.ParseDate(event[0])
			if err != nil {
				slog.Error("Failed to parse date.", "err", err.Error())
				return
			}
			desc := event[1]
			tmpMoney, err := strconv.ParseInt(event[2], 10, 32)
			if err != nil {
				slog.Error(fmt.Sprintf(`failed to parse "%s" as int`, event[2]), "err", err.Error())
				return
			}
			tmpMoney *= -1
			money := int32(tmpMoney)
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
