package console

import (
	"fmt"
	"kakeibodb/internal/model"
)

type TagPresenter struct {
}

func NewTagPresenter() *TagPresenter {
	return &TagPresenter{}
}

func (tp *TagPresenter) Present(tags []*model.TagWithID) {
	fmt.Printf("%s\t%s\n", "ID", "name")
	for _, t := range tags {
		fmt.Printf("%2d\t%v\n", t.GetID(), t.String())
	}
}
