package usecase

import (
	"kakeibodb/internal/model"
	"log"
)

type PatternRepository interface {
	Create(key string) (int64, error)
	Exist(key string) (bool, error)
	Delete(id int64) error
	List() ([]*model.Pattern, error)
}

type PatternTagMapRepository interface {
	Map(patternID int64, tagName string) error
	Unmap(patternID int64, tagName string) error
}

type PatternPresenter interface {
	Present(patterns []*model.Pattern)
}

type PatternUseCase struct {
	patternRepo PatternRepository
}

type PatternPresentUseCase struct {
	PatternUseCase
	patternPresenter PatternPresenter
}

type PatternTagMapUsecase struct {
	ptmRepo PatternTagMapRepository
}

func NewPatternUseCase(patternRepo PatternRepository) *PatternUseCase {
	return &PatternUseCase{
		patternRepo: patternRepo,
	}
}

func NewPatternPresentUseCase(patternRepo PatternRepository, patternPresenter PatternPresenter) *PatternPresentUseCase {
	return &PatternPresentUseCase{
		PatternUseCase:   *NewPatternUseCase(patternRepo),
		patternPresenter: patternPresenter,
	}
}

func NewPatternTagMapUseCase(ptmRepo PatternTagMapRepository) *PatternTagMapUsecase {
	return &PatternTagMapUsecase{
		ptmRepo: ptmRepo,
	}
}

func (pu *PatternUseCase) Create(key string) {
	exist, err := pu.patternRepo.Exist(key)
	if err != nil {
		log.Fatal(err)
	}
	if exist {
		return
	}

	_, err = pu.patternRepo.Create(key)
	if err != nil {
		log.Fatal(err)
	}
}

func (pu *PatternUseCase) Delete(id int64) {
	err := pu.patternRepo.Delete(id)
	if err != nil {
		log.Fatal(err)
	}
}

func (pu *PatternPresentUseCase) List() {
	tags, err := pu.patternRepo.List()
	if err != nil {
		log.Fatal(err)
	}

	pu.patternPresenter.Present(tags)
}

func (ptmu *PatternTagMapUsecase) AddTag(patternID int64, tagNames []string) {
	for _, tagName := range tagNames {
		err := ptmu.ptmRepo.Map(patternID, tagName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (ptmu *PatternTagMapUsecase) RemoveTag(patternID int64, tagName string) {
	err := ptmu.ptmRepo.Unmap(patternID, tagName)
	if err != nil {
		log.Fatal(err)
	}
}
