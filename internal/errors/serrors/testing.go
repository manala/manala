package serrors

import (
	"github.com/stretchr/testify/assert"
)

type Assert struct {
	Type      any
	Message   string
	Arguments []any
	Details   string
	Error     *Assert
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
		s.Equal(assert.Details, _err.ErrorDetails(false))
	} else {
		if assert.Details != "" {
			s.Fail("Error does not contains details")
		}
	}

	// Error
	if _err, ok := err.(interface {
		Unwrap() error
	}); ok {
		__err := _err.Unwrap()
		if __err == nil {
			if assert.Error != nil {
				s.Fail("Error contains a nil error")
			}
		} else {
			if assert.Error == nil {
				s.Fail("Error contains an error")
			} else {
				Equal(s, assert.Error, __err)
			}
		}
	} else {
		if assert.Error != nil {
			s.Fail("Error does not contains an error")
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
