package components

import (
	"errors"
	"manala/internal/serrors"
)

func (s *Suite) TestMessageFromError() {
	tests := []struct {
		test     string
		err      error
		expected *Message
	}{
		{
			test: "Error",
			err:  errors.New("error"),
			expected: &Message{
				Type:    ErrorMessageType,
				Message: "error",
			},
		},
		{
			test: "StructuredError",
			err: serrors.New("structured error").
				WithArguments("foo", "bar").
				WithDetails(`details`).
				WithErrors(
					errors.New("wrapped error"),
					serrors.New("wrapped structured error").
						WithArguments("bar", "baz").
						WithDetails(`wrapped details`),
				),
			expected: &Message{
				Type:    ErrorMessageType,
				Message: "structured error",
				Attributes: []*MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
				Details: `details`,
				Messages: []*Message{
					{
						Type:    ErrorMessageType,
						Message: "wrapped error",
					},
					{
						Type:    ErrorMessageType,
						Message: "wrapped structured error",
						Attributes: []*MessageAttribute{
							{Key: "bar", Value: "baz"},
						},
						Details: `wrapped details`,
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			message := MessageFromError(test.err, false)

			s.Equal(test.expected, message)
		})
	}
}
