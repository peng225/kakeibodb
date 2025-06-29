package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/query/query"
)

type TagRepository struct {
	q *query.Queries
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{
		q: query.New(db),
	}
}

func (tr *TagRepository) Create(tag model.Tag) (int64, error) {
	ctx := context.Background()
	res, err := tr.q.CreateTag(ctx, sql.NullString{
		String: tag.String(),
		Valid:  true,
	})
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (tr *TagRepository) Exist(tag model.Tag) (bool, error) {
	ctx := context.Background()
	_, err := tr.q.GetTag(ctx, sql.NullString{
		String: tag.String(),
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (tr *TagRepository) Delete(id int32) error {
	ctx := context.Background()
	err := tr.q.DeleteTagByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag by ID: %w", err)
	}
	return nil
}
