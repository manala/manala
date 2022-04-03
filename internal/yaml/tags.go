package yaml

import (
	"regexp"
)

var docTagRegex, _ = regexp.Compile(
	`(?m)` +
		`(\s+)` +
		`|([ \t]*#+[ \t]*@(?P<Tag>[a-zA-Z][\w-]*)[ \t]+)` +
		`|([ \t]*#+[ \t]*)` +
		`|(?P<String>.+$)`,
)
var docTagRegexGroups = docTagRegex.SubexpNames()

type DocTag struct {
	Name  string
	Value string
}

type DocTags []*DocTag

func (tags *DocTags) Filter(name string) *DocTags {
	_tags := DocTags{}
	for _, tag := range *tags {
		if tag.Name != name {
			continue
		}
		_tags = append(_tags, tag)
	}
	return &_tags
}
