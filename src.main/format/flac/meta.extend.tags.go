package flac

import (
	"dadp.flactool/types"
	"strconv"
	"strings"
)

type MetaBlockTags struct {
	data     types.SSMap
	matchers map[string]func(rawData string, pattern []string) bool
	children map[string]*MetaBlockTags
}

func NewMetaBlockTags() *MetaBlockTags {
	m := &MetaBlockTags{}

	m.data = types.SSMap{}
	m.matchers = map[string]func(rawData string, pattern []string) bool{}
	m.children = map[string]*MetaBlockTags{}

	return m
}

func (tags *MetaBlockTags) getChild(fullTag string) (*MetaBlockTags, string) {
	for childName, childNode := range tags.children {
		if strings.HasPrefix(fullTag, childName+".") {
			return childNode.getChild(strings.TrimPrefix(fullTag, childName+"."))
		}
	}
	return tags, fullTag
}

func (tags *MetaBlockTags) Get(tag string) string {
	tcNode, tcName := tags.getChild(tag)
	return tcNode.get(tcName)
}

func (tags *MetaBlockTags) get(tag string) string {
	return tags.data.Get(tag)
}

func (tags *MetaBlockTags) Set(tag string, data string, matcher func(rawData string, pattern []string) bool) {
	tcNode, tcName := tags.getChild(tag)
	tcNode.set(tcName, data, matcher)
}

func (tags *MetaBlockTags) set(tag string, data string, matcher func(rawData string, pattern []string) bool) {
	tags.data.Set(tag, data)
	if !tags.data.Has(tag) {
		return
	}

	if matcher != nil {
		tags.matchers[tag] = matcher
	}
}

func (tags *MetaBlockTags) SetChild(prefix string, childNode *MetaBlockTags) {
	if prefix == "" {
		return
	}
	if strings.Index(prefix, ".") > -1 {
		tcNode, tcName := tags.getChild(prefix)
		tcNode.SetChild(tcName, childNode)
	} else {
		tags.setChild(prefix, childNode)
	}
}

func (tags *MetaBlockTags) setChild(prefix string, childNode *MetaBlockTags) {
	tags.children[prefix] = childNode
}

func (tags *MetaBlockTags) Match(tag string, pattern []string) bool {
	tcNode, tcName := tags.getChild(tag)
	return tcNode.match(tcName, pattern)
}

func (tags *MetaBlockTags) match(tag string, pattern []string) bool {
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

func TagMatcher_Bool(rawData string, pattern []string) bool {
	if len(pattern) == 0 {
		return false
	}

	if rawBool, err := strconv.ParseBool(rawData); err != nil {
		return false
	} else if inputBool, err := strconv.ParseBool(rawData); err != nil {
		return false
	} else {
		return rawBool == inputBool
	}
}

func TagMatcher_MetaType(rawData string, pattern []string) bool {
	if len(pattern) == 0 {
		return false
	}

	rawType := strings.TrimSpace(strings.ToLower(rawData))
	inputType := strings.TrimSpace(strings.ToLower(pattern[0]))
	if rawType == inputType {
		return true
	}

	if typeNum, err := strconv.ParseInt(rawType, 10, 32); err == nil {
		rawType = MetaBlockType(typeNum).String()
	}

	if typeNum, err := strconv.ParseInt(inputType, 10, 32); err == nil {
		inputType = MetaBlockType(typeNum).String()
	}

	if rawType == inputType {
		return true
	}

	return false
}
