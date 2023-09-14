package charm

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
			test:    "Empty",
			message: &components.Message{},
			expected: `
			`,
		},
		{
			test: "NoAttributesAndNoDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
			},
			expected: `
				 • message
			`,
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
			expected: `
				 • message
				   foo=bar
			`,
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
			expected: `
				 • message
				   foo=bar

				   details
			`,
		},
		{
			test: "NoAttributesAndDetails",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message",
				Details: "details",
			},
			expected: `
				 • message

				   details
			`,
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
			expected: `
				 • Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				   foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar foo=bar

				   Suspendisse nec sem ligula. Nunc ut quam eros. Interdum et malesuada fames ac ante ipsum primis in faucibus. Donec erat augue, porta et risus non, tempus convallis velit. Quisque sed ligula pharetra, dignissim est ac, pulvinar est. Sed et sapien auctor ipsum faucibus auctor. Etiam ut faucibus enim. In non nibh viverra massa consequat porttitor. Fusce rutrum neque a justo imperdiet lacinia. Vivamus ex felis, ultrices quis diam in, varius suscipit velit. Suspendisse feugiat ante enim, vitae fringilla neque maximus non.
			`,
		},
		{
			test: "Wrapped",
			message: &components.Message{
				Type:    components.InfoMessageType,
				Message: "message 1",
				Attributes: []*components.MessageAttribute{
					{Key: "foo", Value: "bar"},
				},
				Details: "details 1",
				Messages: []*components.Message{
					{
						Type:    components.InfoMessageType,
						Message: "message 2",
						Attributes: []*components.MessageAttribute{
							{Key: "foo", Value: "bar"},
						},
						Details: "details 2",
						Messages: []*components.Message{
							{
								Type:    components.InfoMessageType,
								Message: "message 3",
								Attributes: []*components.MessageAttribute{
									{Key: "foo", Value: "bar"},
								},
								Details: "details 3",
							},
							{
								Type:    components.InfoMessageType,
								Message: "message 4",
								Attributes: []*components.MessageAttribute{
									{Key: "foo", Value: "bar"},
								},
								Details: "details 4",
							},
						},
					},
				},
			},
			expected: `
				 • message 1
				   foo=bar

				   details 1
				   • message 2
				     foo=bar

				     details 2
				     • message 3
				       foo=bar

				       details 3
				     • message 4
				       foo=bar

				       details 4
			`,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			adapter := New(nil, out, err)

			adapter.Message(test.message)

			s.Empty(out)
			heredoc.Equal(s.Assert(), test.expected, err.String())
		})
	}
}
