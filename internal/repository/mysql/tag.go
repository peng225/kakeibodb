package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql/sqlc/query"
)

type TagRepository struct {
	q *query.Queries
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{
		q: query.New(db),
	}
}

func (tr *TagRepository) Create(ctx context.Context, tagName string) (int64, error) {
	tag, err := tr.get(ctx, tagName)
	if err == nil {
		return tag.GetID(), nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	res, err := tr.q.CreateTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (tr *TagRepository) get(ctx context.Context, tagName string) (*model.Tag, error) {
	ret, err := tr.q.GetTag(ctx, sql.NullString{
		String: tagName,
		Valid:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return model.NewTag(ret.ID, ret.Name.String), nil
}

func (tr *TagRepository) Delete(ctx context.Context, id int64) error {
	err := tr.q.DeleteTagByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag by ID: %w", err)
	}
	return nil
}

func (tr *TagRepository) List(ctx context.Context) ([]*model.Tag, error) {
	res, err := tr.q.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	tags := make([]*model.Tag, len(res))
	for i, tag := range res {
		tags[i] = model.NewTag(tag.ID, tag.Name.String)
	}
	return tags, nil
}
