package model

import "slices"

type Pattern struct {
	key  string
	tags []Tag
}

func NewPattern(key string, tags []Tag) *Pattern {
	return &Pattern{
		key:  key,
		tags: tags,
	}
}

func (p *Pattern) AddTag(tags []Tag) {
	for _, tag := range tags {
		if !slices.Contains(p.tags, tag) {
			p.tags = append(p.tags, tag)
		}
	}
}
