package model

import "slices"

type Pattern struct {
	id   int64
	key  string
	tags []Tag
}

func NewPattern(id int64, key string, tags []Tag) *Pattern {
	return &Pattern{
		id:   id,
		key:  key,
		tags: tags,
	}
}

func (p *Pattern) GetKey() string {
	return p.key
}

func (p *Pattern) GetTags() []Tag {
	ret := make([]Tag, len(p.tags))
	copy(ret, p.tags)
	return ret
}

func (p *Pattern) AddTag(tag Tag) {
	if !slices.Contains(p.tags, tag) {
		p.tags = append(p.tags, tag)
	}
}

func (p *Pattern) GetID() int64 {
	return p.id
}
