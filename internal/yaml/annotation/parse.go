package annotation

// Parse parses the given source into annotations.
func Parse(src string) (Annotations, error) {
	var annotations Annotations
	var current *Annotation

	scanner := NewScanner(src)

	for {
		token := scanner.Scan()

		switch token.Kind {
		case TokenName:
			current = &Annotation{nameToken: token}
			annotations = append(annotations, current)
		case TokenText:
			if current == nil {
				continue
			}
			current.valueTokens = append(current.valueTokens, token)
		case TokenUnknown:
			continue
		case TokenEOF:
			return annotations, nil
		}
	}
}
