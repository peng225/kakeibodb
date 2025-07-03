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
		tagsToPrint := "NONE"
		if len(p.GetTags()) != 0 {
			strTags := make([]string, len(p.GetTags()))
			for i, tag := range p.GetTags() {
				strTags[i] = tag.String()
			}
			tagsToPrint = strings.Join(strTags, ",")
		}
		fmt.Printf("%v\t%-32s\t%s\n",
			p.GetID(), p.GetKey(), tagsToPrint)
	}
}
