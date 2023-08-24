package lipgloss

import (
	"bytes"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/components"
)

func (s *Suite) TestMessage() {
	tests := []struct {
		test     string
		message  *components.Message
		expected string
	}{
		{
			test:     "Empty",
			message:  &components.Message{},
			expected: "",
		},
		{
			test: "NoAttributesAndNoDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
			},
			expected: heredoc.Doc(`
				  • message
			`),
		},
		{
			test: "AttributesAndNoDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
			},
			expected: heredoc.Doc(`
				  • message                            foo=bar
			`),
		},
		{
			test: "AttributesAndDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
				Details: "details",
			},
			expected: heredoc.Doc(`
				  • message                            foo=bar

				    details
			`),
		},
		{
			test: "NoAttributesAndDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
				Details: "details",
			},
			expected: heredoc.Doc(`
				  • message

				    details
			`),
		},
		{
			test: "Large",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
					{Key: "foo", Value: "bar"},
				},
				Details: "Suspendisse nec sem ligula. Nunc ut quam eros. Interdum et malesuada fames ac ante ipsum primis in faucibus. Donec erat augue, porta et risus non, tempus convallis velit. Quisque sed ligula pharetra, dignissim est ac, pulvinar est. Sed et sapien auctor ipsum faucibus auctor. Etiam ut faucibus enim. In non nibh viverra massa consequat porttitor. Fusce rutrum neque a justo imperdiet lacinia. Vivamus ex felis, ultrices quis diam in, varius suscipit velit. Suspendisse feugiat ante enim, vitae fringilla neque maximus non.",
			},
			expected: heredoc.Doc(`
				  • Lorem ipsum dolor sit amet,        foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar
				    consectetur adipiscing elit.

				    Suspendisse nec sem ligula. Nunc ut quam eros. Interdum et malesuada fames ac ante ipsum primis in faucibus. Donec erat augue, porta et risus non, tempus convallis velit. Quisque sed ligula pharetra, dignissim est ac, pulvinar est. Sed et sapien auctor ipsum faucibus auctor. Etiam ut faucibus enim. In non nibh viverra massa consequat porttitor. Fusce rutrum neque a justo imperdiet lacinia. Vivamus ex felis, ultrices quis diam in, varius suscipit velit. Suspendisse feugiat ante enim, vitae fringilla neque maximus non.
			`),
		},
		{
			test: "Wrapped",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
				Details: "details",
				Messages: []*components.Message{
					{
						Type:    components.InfoMessageType,
						Message: "message",
						Attributes: []*components.MessageAttribute{
							{Key: "foo", Value: "bar"},
						},
						Details: "details",
						Messages: []*components.Message{
							{
								Type:    components.InfoMessageType,
								Message: "message",
								Attributes: []*components.MessageAttribute{
									{Key: "foo", Value: "bar"},
								},
								Details: "details",
							},
							{
								Type:    components.InfoMessageType,
								Message: "message",
								Attributes: []*components.MessageAttribute{
									{Key: "foo", Value: "bar"},
								},
								Details: "details",
							},
						},
					},
				},
			},
			expected: heredoc.Doc(`
				  • message                            foo=bar

				    details
				    • message                            foo=bar

				      details
				      • message                            foo=bar

				        details
				      • message                            foo=bar

				        details
			`),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			output := New(out, err)

			output.Message(test.message)

			s.Empty(out)
			s.Equal(test.expected, err.String())
		})
	}
}
