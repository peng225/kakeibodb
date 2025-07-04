package usecase

import (
	"kakeibodb/internal/model"
	"log"
)

type TagRepository interface {
	Create(tag model.Tag) (int64, error)
	Delete(id int64) error
	List() ([]*model.TagWithID, error)
}

type TagPresenter interface {
	Present(tags []*model.TagWithID)
}

type TagUseCase struct {
	tagRepo TagRepository
}

type TagPresentUseCase struct {
	TagUseCase
	tagPresenter TagPresenter
}

func NewTagUseCase(tagRepo TagRepository) *TagUseCase {
	return &TagUseCase{
		tagRepo: tagRepo,
	}
}

func NewTagPresentUseCase(tagRepo TagRepository, tagPresenter TagPresenter) *TagPresentUseCase {
	return &TagPresentUseCase{
		TagUseCase:   *NewTagUseCase(tagRepo),
		tagPresenter: tagPresenter,
	}
}

func (tu *TagUseCase) Create(tag model.Tag) {
	_, err := tu.tagRepo.Create(tag)
	if err != nil {
		log.Fatal(err)
	}
}

func (tu *TagUseCase) Delete(id int64) {
	err := tu.tagRepo.Delete(id)
	if err != nil {
		log.Fatal(err)
	}
}

func (tu *TagPresentUseCase) List() {
	tags, err := tu.tagRepo.List()
	if err != nil {
		log.Fatal(err)
	}

	tu.tagPresenter.Present(tags)
}
