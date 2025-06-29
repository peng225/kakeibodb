package usecase

import (
	"kakeibodb/internal/model"
	"log"
)

type TagRepository interface {
	Create(tag model.Tag) (int64, error)
	Exist(tag model.Tag) (bool, error)
	Delete(id int32) error
}

type TagUseCase struct {
	tagRepo TagRepository
}

func NewTagUseCase(tagRepo TagRepository) *TagUseCase {
	return &TagUseCase{
		tagRepo: tagRepo,
	}
}

func (tu *TagUseCase) Create(tag model.Tag) {
	exist, err := tu.tagRepo.Exist(tag)
	if err != nil {
		log.Fatal(err)
	}
	// FIXME: Need to resolve TOCTOU by using transaction.
	if exist {
		return
	}

	_, err = tu.tagRepo.Create(tag)
	if err != nil {
		log.Fatal(err)
	}
}

func (tu *TagUseCase) Delete(id int32) {
	err := tu.tagRepo.Delete(id)
	if err != nil {
		log.Fatal(err)
	}
}
