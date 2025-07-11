package usecase_test

import (
	"context"
	"kakeibodb/internal/repository/mysql/fake"
	"kakeibodb/internal/usecase"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplit_SingleEventID(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	tx := fake.NewFakeTransaction()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "fruit",
	})
	require.NoError(t, err)
	eventUC := usecase.NewEventUseCase(fakeEventRepo, tx)
	err = eventUC.Split(ctx, []int64{0}, "",
		time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local), -30,
		"apple")
	require.NoError(t, err)
	e0, err := fakeEventRepo.GetWithoutTags(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(-70), e0.GetMoney())
	e1, err := fakeEventRepo.GetWithoutTags(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(-30), e1.GetMoney())
}

func TestSplit_SingleEventID_SameMoney(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	tx := fake.NewFakeTransaction()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "fruit",
	})
	require.NoError(t, err)
	eventUC := usecase.NewEventUseCase(fakeEventRepo, tx)
	err = eventUC.Split(ctx, []int64{0}, "",
		time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local), -100,
		"apple")
	require.NoError(t, err)
	_, err = fakeEventRepo.GetWithoutTags(ctx, 0)
	require.Error(t, err)
	e1, err := fakeEventRepo.GetWithoutTags(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(-100), e1.GetMoney())
}

func TestSplit_TwoEventIDs(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	tx := fake.NewFakeTransaction()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "fruit1",
	})
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local),
		Money: int32(-200),
		Desc:  "fruit2",
	})
	require.NoError(t, err)
	eventUC := usecase.NewEventUseCase(fakeEventRepo, tx)
	err = eventUC.Split(ctx, []int64{0, 1}, "",
		time.Date(2025, 6, 7, 0, 0, 0, 0, time.Local), -160,
		"apple")
	require.NoError(t, err)
	_, err = fakeEventRepo.GetWithoutTags(ctx, 0)
	require.Error(t, err)
	e1, err := fakeEventRepo.GetWithoutTags(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(-140), e1.GetMoney())
	e2, err := fakeEventRepo.GetWithoutTags(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, int32(-160), e2.GetMoney())
}

func TestSplit_SingleEvent_AutoDetect(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	tx := fake.NewFakeTransaction()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "fruit",
	})
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 0, "fruit")
	require.NoError(t, err)
	eventUC := usecase.NewEventUseCase(fakeEventRepo, tx)

	err = os.Setenv("KAKEIBODB_SPLIT_BASE_TAG_NAME", "fruit")
	require.NoError(t, err)
	t.Cleanup(func() { os.Unsetenv("KAKEIBODB_SPLIT_BASE_TAG_NAME") })
	err = eventUC.Split(ctx, nil, "fruit",
		time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local), -30,
		"apple")
	require.NoError(t, err)
	e0, err := fakeEventRepo.GetWithoutTags(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(-70), e0.GetMoney())
	e1, err := fakeEventRepo.GetWithoutTags(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int32(-30), e1.GetMoney())
}

func TestSplit_TwoEventIDs_AutoDetect(t *testing.T) {
	fakeEventRepo := fake.NewEventFakeRepository()
	tx := fake.NewFakeTransaction()
	ctx := context.Background()
	_, err := fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local),
		Money: int32(-100),
		Desc:  "fruit1",
	})
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 0, "fruit")
	require.NoError(t, err)
	_, err = fakeEventRepo.Create(ctx, &usecase.EventCreateRequest{
		Date:  time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local),
		Money: int32(-200),
		Desc:  "fruit2",
	})
	require.NoError(t, err)
	err = fakeEventRepo.AddTag(ctx, 1, "fruit")
	require.NoError(t, err)
	eventUC := usecase.NewEventUseCase(fakeEventRepo, tx)
	err = os.Setenv("KAKEIBODB_SPLIT_BASE_TAG_NAME", "fruit")
	require.NoError(t, err)
	t.Cleanup(func() { os.Unsetenv("KAKEIBODB_SPLIT_BASE_TAG_NAME") })
	err = eventUC.Split(ctx, nil, "fruit",
		time.Date(2025, 6, 7, 0, 0, 0, 0, time.Local), -260,
		"apple")
	require.NoError(t, err)
	e0, err := fakeEventRepo.GetWithoutTags(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(-40), e0.GetMoney())
	_, err = fakeEventRepo.GetWithoutTags(ctx, 1)
	require.Error(t, err)
	e2, err := fakeEventRepo.GetWithoutTags(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, int32(-260), e2.GetMoney())
}
