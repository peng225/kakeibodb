package model

import "slices"

type Pattern struct {
	id       int64
	key      string
	tagNames []string
}

func NewPattern(id int64, key string, tagNames []string) *Pattern {
	return &Pattern{
		id:       id,
		key:      key,
		tagNames: tagNames,
	}
}

func (p *Pattern) GetKey() string {
	return p.key
}

func (p *Pattern) GetTagNames() []string {
	ret := make([]string, len(p.tagNames))
	copy(ret, p.tagNames)
	return ret
}

func (p *Pattern) AddTag(tagName string) {
	if !slices.Contains(p.tagNames, tagName) {
		p.tagNames = append(p.tagNames, tagName)
	}
}

func (p *Pattern) GetID() int64 {
	return p.id
}
