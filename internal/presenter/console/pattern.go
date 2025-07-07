package console

import (
	"fmt"
	"kakeibodb/internal/model"
	"strings"
)

type PatternPresenter struct {
}

func NewPatternPresenter() *PatternPresenter {
	return &PatternPresenter{}
}

func (pp *PatternPresenter) Present(patterns []*model.Pattern) {
	fmt.Printf("%s\t%s\t%s\n", "ID", "key                             ", "tags")
	for _, p := range patterns {
		tagsToPrint := model.EmptyTagName
		if len(p.GetTagNames()) != 0 {
			tagsToPrint = strings.Join(p.GetTagNames(), ",")
		}
		fmt.Printf("%v\t%-32s\t%s\n",
			p.GetID(), p.GetKey(), tagsToPrint)
	}
}
