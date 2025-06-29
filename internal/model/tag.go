package model

type Tag string

type TagWithID struct {
	tag Tag
	id  int32
}

func (t *Tag) String() string {
	return string(*t)
}

func NewTagWithID(id int32, tag Tag) *TagWithID {
	return &TagWithID{
		tag: tag,
		id:  id,
	}
}

func (t *TagWithID) GetID() int32 {
	return t.id
}

func (t *TagWithID) String() string {
	return t.tag.String()
}
