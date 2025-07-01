package usecase

import (
	"kakeibodb/internal/model"
	"log"
)

type PatternRepository interface {
	Create(key string) (int64, error)
	Exist(key string) (bool, error)
	Delete(id int64) error
}

type PatternPresenter interface {
	Present(patterns []*model.PatternWithID)
}

type PatternUseCase struct {
	patternRepo PatternRepository
}

type PatternPresentUseCase struct {
	PatternUseCase
	patternPresenter PatternPresenter
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
