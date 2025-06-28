package model

type Tag string

func (t *Tag) String() string {
	return string(*t)
}
