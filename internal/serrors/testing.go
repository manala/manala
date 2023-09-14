package serrors

import (
	"github.com/stretchr/testify/assert"
	"manala/internal/testing/heredoc"
)

type Assert struct {
	Type      any
	Message   string
	Arguments []any
	Details   string
	Errors    []*Assert
}

func Equal(s *assert.Assertions, assert *Assert, err error) {
	s.IsType(assert.Type, err)
	s.EqualError(err, assert.Message)

	// Arguments
	if _err, ok := err.(ErrorArguments); ok {
		s.Equal(assert.Arguments, _err.ErrorArguments())
	} else {
		if assert.Arguments != nil {
			s.Fail("Error does not contains arguments")
		}
	}

	// Details
	if _err, ok := err.(ErrorDetails); ok {
		heredoc.Equal(s, assert.Details, _err.ErrorDetails(false))
	} else {
		if assert.Details != "" {
			s.Fail("Error does not contains details")
		}
	}

	// Errors
	if _err, ok := err.(interface {
		Unwrap() []error
	}); ok {
		__errs := _err.Unwrap()
		if __errs == nil {
			if assert.Errors != nil {
				s.Fail("Error contains nil errors")
			}
		} else {
			if assert.Errors == nil {
				s.Fail("Error contains errors")
			} else {
				s.Len(__errs, len(assert.Errors), "Incorrect error's errors length")
				for i, _assert := range assert.Errors {
					Equal(s, _assert, __errs[i])
				}
			}
		}
	} else {
		if assert.Errors != nil {
			s.Fail("Error does not contains errors")
		}
	}
}
