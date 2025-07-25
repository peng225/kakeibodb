package usecase_test

import (
	"context"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/fake"
	"kakeibodb/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMoneySumGroupedByTagName_NoTag(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	analysisUC := usecase.NewAnalysisUseCase(fakeEventRepo, nil)
	from := time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)
	to := time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local)
	ret, err := analysisUC.GetMoneySumGroupedByTagName(ctx, from, to)
	require.NoError(t, err)
	require.Len(t, ret, 1)
	assert.Equal(t, int32(-100), ret[model.EmptyTagName])
}

func TestGetMoneySumGroupedByTagName_BoundaryCheck(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 5, 31, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 30, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 0, "fruit")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 1, "fruit")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 2, "fruit")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 3, "fruit")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 4, "fruit")
	require.NoError(t, err)
	analysisUC := usecase.NewAnalysisUseCase(fakeEventRepo, nil)
	from := time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)
	to := time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local)
	ret, err := analysisUC.GetMoneySumGroupedByTagName(ctx, from, to)
	require.NoError(t, err)
	require.Len(t, ret, 1)
	assert.Equal(t, int32(-300), ret["fruit"])
}

func TestGetMoneySumGroupedByTagName_MultipleTags(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "apple",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local),
		Money: int32(-50),
		Desc:  "pencil",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local),
		Money: int32(-200),
		Desc:  "orange",
	})
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 0, "fruit")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 1, "stationary")
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 2, "fruit")
	require.NoError(t, err)
	analysisUC := usecase.NewAnalysisUseCase(fakeEventRepo, nil)
	from := time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)
	to := time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local)
	ret, err := analysisUC.GetMoneySumGroupedByTagName(ctx, from, to)
	require.NoError(t, err)
	require.Len(t, ret, 2)
	assert.Equal(t, int32(-300), ret["fruit"])
	assert.Equal(t, int32(-50), ret["stationary"])
}

func TestGetHighlyRankedTagNames_Empty(t *testing.T) {
	msGroupedByTagNameForEveryWindow := make([](map[string]int32), 0)
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 10)
	require.Empty(t, tagNames)
}

func TestGetHighlyRankedTagNames_SmallerNumberOfItemsThanTopVar(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 10)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "fruit", tagNames[0])
	assert.Equal(t, "stationary", tagNames[1])
}

func TestGetHighlyRankedTagNames_RankIsUpdatedByTheLastItem(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
			"toy":        -200,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 2)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "fruit", tagNames[0])
	assert.Equal(t, "toy", tagNames[1])
}

func TestGetHighlyRankedTagNames_RankIsNotUpdatedByTheLastItem(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
			"toy":        -80,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 2)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "fruit", tagNames[0])
	assert.Equal(t, "stationary", tagNames[1])
}

func TestGetHighlyRankedTagNames_UpdatedNotRankedItemsByLaterWindow(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
			"toy":        -200,
		},
		{
			"stationary": -400,
			"sports":     -500,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 2)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "sports", tagNames[0])
	assert.Equal(t, "stationary", tagNames[1])
}

func TestGetHighlyRankedTagNames_UpdatedAlreadyRankedItemByLaterWindow2(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
			"toy":        -200,
		},
		{
			"toy": -400,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 2)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "toy", tagNames[0])
	assert.Equal(t, "fruit", tagNames[1])
}

func TestGetHighlyRankedTagNames_NotUpdatedByLaterWindow(t *testing.T) {
	msGroupedByTagNameForEveryWindow := [](map[string]int32){
		map[string]int32{
			"stationary": -100,
			"fruit":      -300,
			"toy":        -200,
		},
		{
			"toy":    -80,
			"sports": -50,
		},
	}
	tagNames := usecase.GetHighlyRankedTagNames(msGroupedByTagNameForEveryWindow, 2)
	require.Len(t, tagNames, 2)
	assert.Equal(t, "fruit", tagNames[0])
	assert.Equal(t, "toy", tagNames[1])
}
