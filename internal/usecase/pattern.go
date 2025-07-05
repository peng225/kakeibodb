package usecase

import (
	"fmt"
	"kakeibodb/internal/model"
	"log/slog"
)

type PatternRepository interface {
	Create(key string) (int64, error)
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

func (pu *PatternUseCase) Create(key string) error {
	id, err := pu.patternRepo.Create(key)
	if err != nil {
		return fmt.Errorf("failed to create pattern: %w", err)
	}
	slog.Info("Created a pattern.", "ID", id)
	return nil
}

func (pu *PatternUseCase) Delete(id int64) error {
	err := pu.patternRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete pattern: %w", err)
	}
	return nil
}

func (pu *PatternPresentUseCase) List() error {
	tags, err := pu.patternRepo.List()
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	pu.patternPresenter.Present(tags)
	return nil
}

func (ptmu *PatternTagMapUsecase) AddTag(patternID int64, tagNames []string) error {
	for _, tagName := range tagNames {
		err := ptmu.ptmRepo.Map(patternID, tagName)
		if err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}
	}
	return nil
}

func (ptmu *PatternTagMapUsecase) RemoveTag(patternID int64, tagName string) error {
	err := ptmu.ptmRepo.Unmap(patternID, tagName)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}
	return nil
}
