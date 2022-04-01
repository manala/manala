package doc

import (
	"regexp"
)

var regex, _ = regexp.Compile(
	`(?m)` +
		`(\s+)` +
		`|([ \t]*#+[ \t]*@(?P<Tag>[a-zA-Z][\w-]*)[ \t]+)` +
		`|([ \t]*#+[ \t]*)` +
		`|(?P<String>.+$)`,
)
var regexGroups = regex.SubexpNames()

type Tag struct {
	Name  string
	Value string
}

type TagList struct {
	tags []*Tag
}

func (l *TagList) Add(tag *Tag) {
	l.tags = append(l.tags, tag)
}

func (l *TagList) All() []*Tag {
	return l.tags
}

func (l *TagList) Filter(name string) []*Tag {
	tags := make([]*Tag, 0)
	for _, tag := range l.tags {
		if tag.Name != name {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func ParseCommentTags(comment string) TagList {
	var list TagList

	// Lexer
	tokens := make([][2]string, 0)
	for _, submatch := range regex.FindAllStringSubmatch(comment, -1) {
		for i, match := range submatch {
			if i != 0 && match != "" {
				group := regexGroups[i]
				if group != "" {
					tokens = append(tokens, [2]string{group, match})
				}
			}
		}
	}
	// Parser
	tag := &Tag{}
	for _, token := range tokens {
		switch token[0] {
		case "Tag":
			if tag.Name != "" {
				if tag.Value != "" {
					list.Add(tag)
				}
				tag = &Tag{}
			}
			tag.Name = token[1]
		case "String":
			if tag.Name != "" {
				if tag.Value == "" {
					tag.Value = token[1]
				} else {
					tag.Value += "\n" + token[1]
				}
			}
		}
	}
	if tag.Name != "" && tag.Value != "" {
		list.Add(tag)
	}

	return list
}
