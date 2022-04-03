package yaml

func ParseComment(comment string, docTags *DocTags) {
	// Lexer
	tokens := make([][2]string, 0)
	for _, submatch := range docTagRegex.FindAllStringSubmatch(comment, -1) {
		for i, match := range submatch {
			if i != 0 && match != "" {
				group := docTagRegexGroups[i]
				if group != "" {
					tokens = append(tokens, [2]string{group, match})
				}
			}
		}
	}
	// Parser
	docTag := &DocTag{}
	for _, token := range tokens {
		switch token[0] {
		case "Tag":
			if docTag.Name != "" {
				if docTag.Value != "" {
					*docTags = append(*docTags, docTag)
				}
				docTag = &DocTag{}
			}
			docTag.Name = token[1]
		case "String":
			if docTag.Name != "" {
				if docTag.Value == "" {
					docTag.Value = token[1]
				} else {
					docTag.Value += "\n" + token[1]
				}
			}
		}
	}
	if docTag.Name != "" && docTag.Value != "" {
		*docTags = append(*docTags, docTag)
	}
}
