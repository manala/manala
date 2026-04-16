package annotation

import (
	"errors"
	"fmt"
)

// Parse parses the given source into a Set of annotations.
func Parse(src string) (*Set, error) {
	set := &Set{}
	var current *Annotation

	scanner := NewScanner(src)
	seen := map[string]bool{}

	for {
		token := scanner.Scan()

		switch token.Kind {
		case TokenName:
			if seen[token.Value] {
				return nil, ErrorAt(
					fmt.Errorf("duplicate annotation @%s", token.Value),
					token,
				)
			}
			seen[token.Value] = true
			current = &Annotation{Name: Name{Token: token}}
			set.annotations = append(set.annotations, current)
		case TokenText:
			if current == nil {
				continue
			}
			current.Value.Tokens = append(current.Value.Tokens, token)
		case TokenUnknown:
			return nil, ErrorAt(
				errors.New("unknown annotation token"),
				token,
			)
		case TokenEOF:
			return set, nil
		}
	}
}
