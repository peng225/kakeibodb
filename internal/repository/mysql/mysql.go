package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/query/query"
)

type MySQLRepository struct {
	q *query.Queries
}

func NewEventRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		q: query.New(db),
	}
}

func (mr *MySQLRepository) Create(event *model.Event) (int64, error) {
	ctx := context.Background()
	res, err := mr.q.CreateEvent(ctx, query.CreateEventParams{
		Dt: sql.NullTime{
			Time:  event.GetDate(),
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: event.GetMoney(),
			Valid: true,
		},
		Description: sql.NullString{
			String: event.GetDesc(),
			Valid:  true,
		},
	})
	if err != nil {
		return 0, err
	}
	// FIXME: tags are ignored. Is the model correct?
	return res.LastInsertId()
}

func (mr *MySQLRepository) Exist(event *model.Event) (bool, error) {
	ctx := context.Background()
	_, err := mr.q.GetEvent(ctx, query.GetEventParams{
		Dt: sql.NullTime{
			Time:  event.GetDate(),
			Valid: true,
		},
		Money: sql.NullInt32{
			Int32: event.GetMoney(),
			Valid: true,
		},
		Description: sql.NullString{
			String: event.GetDesc(),
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (mr *MySQLRepository) Get(id int32) (*model.Event, error) {
	ctx := context.Background()
	res, err := mr.q.GetEventByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}
	return model.NewEvent(res.Dt.Time, res.Money.Int32, res.Description.String, nil), nil
}

func (mr *MySQLRepository) Delete(id int32) error {
	ctx := context.Background()
	err := mr.q.DeleteEventByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete event by ID: %w", err)
	}
	return nil
}
