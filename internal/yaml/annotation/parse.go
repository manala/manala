package annotation

import (
	"errors"
	"fmt"
)

// Parse parses the given source into a Set of annotations.
func Parse(src string) ([]*Annotation, error) {
	var annotations []*Annotation
	var current *Annotation

	scanner := NewScanner(src)
	seen := map[string]bool{}

	for {
		token := scanner.Scan()

		switch token.Kind {
		case TokenName:
			if seen[token.Value] {
				return nil, NewError(
					fmt.Errorf("duplicate @%s annotation", token.Value),
					token,
				)
			}
			seen[token.Value] = true
			current = &Annotation{Name: Name{Token: token}}
			annotations = append(annotations, current)
		case TokenText:
			if current == nil {
				continue
			}
			if current.Body == nil {
				current.Body = &Body{}
			}
			current.Body.Tokens = append(current.Body.Tokens, token)
		case TokenUnknown:
			return nil, NewError(
				errors.New("unknown annotation token"),
				token,
			)
		case TokenEOF:
			return annotations, nil
		}
	}
}
