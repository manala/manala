package components_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/ui/components"

	"github.com/stretchr/testify/suite"
)

type MessageSuite struct{ suite.Suite }

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

func (s *MessageSuite) TestFromError() {
	tests := []struct {
		test     string
		err      error
		expected *components.Message
	}{
		{
			test: "Error",
			err:  errors.New("error"),
			expected: &components.Message{
				Type:    components.ErrorMessageType,
				Message: "error",
			},
		},
		{
			test: "StructuredError",
			err: serrors.New("structured error").
				WithArguments("foo", "bar").
				WithDump(`dump`).
				WithErrors(
					errors.New("wrapped error"),
					serrors.New("wrapped structured error").
						WithArguments("bar", "baz").
						WithDump(`wrapped dump`),
				),
			expected: &components.Message{
				Type:    components.ErrorMessageType,
				Message: "structured error",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
				Dump: `dump`,
				Messages: []*components.Message{
					{
						Type:    components.ErrorMessageType,
						Message: "wrapped error",
					},
					{
						Type:    components.ErrorMessageType,
						Message: "wrapped structured error",
						Attributes: []*components.MessageAttribute{
							{Key: "bar", Value: "baz"},
						},
						Dump: `wrapped dump`,
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			message := components.MessageFromError(test.err, false)

			s.Equal(test.expected, message)
		})
	}
}
