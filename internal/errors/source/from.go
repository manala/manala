package source

import "errors"

// From collects all Positions in the error tree and wraps each in an Error with the given origin.
// If no Position is found, the original error is returned as-is.
// The returned error is an errors.Join of all collected Errors — each one walks its own chain via Unwrap().
func From(err error, origin Origin) error {
	positions := from(err)
	if len(positions) == 0 {
		// no Position anywhere in the tree — return the original error
		return err
	}
	errs := make([]error, len(positions))
	for i, pos := range positions {
		absLine, absCol := pos.Position()
		errs[i] = Error{Origin: origin, Position: pos, Line: absLine, Column: absCol}
	}
	return errors.Join(errs...)
}

// from follows the same traversal logic as errors.As — Unwrap() error is followed
// iteratively, Unwrap() []error fans out recursively into each branch — but instead
// of stopping at the first match it collects every Position found, one per branch.
func from(err error) []Position {
	for err != nil {
		if pos, ok := err.(Position); ok {
			// found one — stop here and collect it
			return []Position{pos}
		}
		switch x := err.(type) {
		case interface{ Unwrap() []error }:
			// fan out: each branch may independently contain a Position
			var result []Position
			for _, child := range x.Unwrap() {
				result = append(result, from(child)...)
			}
			return result
		case interface{ Unwrap() error }:
			// not a Position — keep descending
			err = x.Unwrap()
		default:
			// true leaf with no Position
			return nil
		}
	}
	return nil
}
