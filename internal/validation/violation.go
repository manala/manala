package validation

import "strings"

type Violation struct { //nolint:errname
	error

	location     string
	line, column int
}

func (e Violation) Error() string {
	return e.error.Error()
}

func (e Violation) Location() string {
	return e.location
}

func (e Violation) Position() (int, int) {
	return e.line, e.column
}

func (e Violation) Unwrap() error {
	return e.error
}

type Violations []*Violation //nolint:errname

func (v Violations) Error() string {
	msgs := make([]string, len(v))
	for i, violation := range v {
		msgs[i] = violation.Error()
	}
	return strings.Join(msgs, "\n")
}

func (v Violations) Unwrap() []error {
	errs := make([]error, len(v))
	for i, violation := range v {
		errs[i] = violation
	}
	return errs
}

func (v Violations) First() (*Violation, bool) {
	if len(v) == 0 {
		return nil, false
	}
	return v[0], true
}
