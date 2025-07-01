package model

import "slices"

type Pattern struct {
	key  string
	tags []Tag
}

type PatternWithID struct {
	Pattern
	id int64
}

func NewPattern(key string, tags []Tag) *Pattern {
	return &Pattern{
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

func NewPatternWithID(id int64, key string, tags []Tag) *PatternWithID {
	return &PatternWithID{
		Pattern: *NewPattern(key, tags),
		id:      id,
	}
}

func (p *PatternWithID) GetID() int64 {
	return p.id
}
