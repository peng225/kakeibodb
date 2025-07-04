package model

type Tag struct {
	id   int64
	name string
}

func NewTag(id int64, name string) *Tag {
	return &Tag{
		id:   id,
		name: name,
	}
}

func (t *Tag) GetName() string {
	return t.name
}

func (t *Tag) GetID() int64 {
	return t.id
}
