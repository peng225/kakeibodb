package mysql

import (
	"context"
	"database/sql"
	"errors"
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
