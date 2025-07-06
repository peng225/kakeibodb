package usecase

import (
	"fmt"
	"kakeibodb/internal/model"
)

type TagPresenter interface {
	Present(tags []*model.Tag)
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

func (tu *TagUseCase) Create(tagName string) error {
	if !model.ValidTagName(tagName) {
		return fmt.Errorf("invalid tag name (%s)", tagName)
	}
	_, err := tu.tagRepo.Create(tagName)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (tu *TagUseCase) Delete(id int64) error {
	err := tu.tagRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return err
}

func (tu *TagPresentUseCase) List() error {
	tags, err := tu.tagRepo.List()
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}

	tu.tagPresenter.Present(tags)
	return nil
}
