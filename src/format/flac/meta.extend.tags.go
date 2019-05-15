package flac

import (
	"p20190417/types"
	"strings"
)

type MetaBlockTags struct {
	data     types.SSMap
	matchers map[string]func(rawData string, pattern []string) bool
}

func (tags *MetaBlockTags) String() string {
	separator := ", "
	result := types.NewBuffer()

	for tag, tagContent := range tags.data.Dump() {
		_, _ = result.WriteStrings(tag, "=", tagContent, separator)
	}

	return strings.TrimSuffix(result.String(), separator)
}

func NewMetaBlockTags() *MetaBlockTags {
	m := &MetaBlockTags{}

	m.data = types.SSMap{}
	m.matchers = map[string]func(rawData string, pattern []string) bool{}

	return m
}

func (tags *MetaBlockTags) Set(tag string, data string, matcher func(rawData string, pattern []string) bool) {
	tags.SetData(tag, data)
	tags.SetMatcher(tag, matcher)
}

func (tags *MetaBlockTags) SetData(tag string, data string) {
	tags.data.Set(tag, data)
}

func (tags *MetaBlockTags) SetMatcher(tag string, matcher func(rawData string, pattern []string) bool) {
	if !tags.data.Has(tag) {
		return
	}

	if matcher != nil {
		tags.matchers[tag] = matcher
	}
}

func (tags *MetaBlockTags) Match(tag string, pattern []string) bool {
	if tags.data.Has(tag) {
		tagContent := tags.data.Get(tag)
		if matcher, matcherExist := tags.matchers[tag]; matcherExist && matcher != nil {
			return matcher(tagContent, pattern)
		} else {
			return len(pattern) == 1 && tagContent == pattern[0]
		}
	} else {
		return false
	}
}

func (tags *MetaBlockTags) Matches(patterns map[string][]string) bool {
	for tag, pattern := range patterns {
		if !tags.Match(tag, pattern) {
			return false
		}
	}
	return true
}
