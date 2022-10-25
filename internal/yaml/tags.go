package yaml

import "regexp"

func ParseCommentTags(comment string, tags *Tags) {
	// Lexer
	tokens := make([][2]string, 0)
	for _, submatch := range tagRegex.FindAllStringSubmatch(comment, -1) {
		for i, match := range submatch {
			if i != 0 && match != "" {
				group := tagRegexGroups[i]
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
					*tags = append(*tags, tag)
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
		*tags = append(*tags, tag)
	}
}

var tagRegex, _ = regexp.Compile(
	`(?m)` +
		`(\s+)` +
		`|([ \t]*#+[ \t]*@(?P<Tag>[a-zA-Z][\w-]*)[ \t]+)` +
		`|([ \t]*#+[ \t]*)` +
		`|(?P<String>.+$)`,
)
var tagRegexGroups = tagRegex.SubexpNames()

type Tag struct {
	Name  string
	Value string
}

type Tags []*Tag

func (tags *Tags) Filter(name string) *Tags {
	_tags := Tags{}
	for _, tag := range *tags {
		if tag.Name != name {
			continue
		}
		_tags = append(_tags, tag)
	}
	return &_tags
}
