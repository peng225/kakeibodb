package usecase

import (
	"context"
	"fmt"
	"kakeibodb/internal/model"
	"log/slog"
)

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

func (pu *PatternUseCase) Create(ctx context.Context, key string) error {
	id, err := pu.patternRepo.Create(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to create pattern: %w", err)
	}
	slog.Info("Created a pattern.", "ID", id)
	return nil
}

func (pu *PatternUseCase) Delete(ctx context.Context, id int64) error {
	err := pu.patternRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pattern: %w", err)
	}
	return nil
}

func (pu *PatternPresentUseCase) List(ctx context.Context) error {
	tags, err := pu.patternRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list patterns: %w", err)
	}

	pu.patternPresenter.Present(tags)
	return nil
}

func (ptmu *PatternTagMapUsecase) AddTag(ctx context.Context, patternID int64, tagNames []string) error {
	for _, tagName := range tagNames {
		err := ptmu.ptmRepo.Map(ctx, patternID, tagName)
		if err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}
	}
	return nil
}

func (ptmu *PatternTagMapUsecase) RemoveTag(ctx context.Context, patternID int64, tagName string) error {
	err := ptmu.ptmRepo.Unmap(ctx, patternID, tagName)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}
	return nil
}
