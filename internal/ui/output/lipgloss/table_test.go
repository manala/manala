package lipgloss

import (
	"bytes"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/components"
)

func (s *Suite) TestTable() {
	tests := []struct {
		test     string
		table    *components.Table
		expected string
	}{
		{
			test:     "Empty",
			table:    &components.Table{},
			expected: "",
		},
		{
			test: "single",
			table: &components.Table{
				Rows: []*components.TableRow{
					{
						Primary:   "primary",
						Secondary: "secondary",
					},
				},
			},
			expected: heredoc.Doc(`
				primary  secondary
			`),
		},
		{
			test: "Multiple",
			table: &components.Table{
				Rows: []*components.TableRow{
					{
						Primary:   "primary 1",
						Secondary: "secondary 1",
					},
					{
						Primary:   "primary 2",
						Secondary: "secondary 2",
					},
				},
			},
			expected: heredoc.Doc(`
				primary 1  secondary 1
				primary 2  secondary 2
			`),
		},
		{
			test: "Alignment",
			table: &components.Table{
				Rows: []*components.TableRow{
					{
						Primary:   "foo",
						Secondary: "bar",
					},
					{
						Primary:   "foofoo",
						Secondary: "barbar",
					},
				},
			},
			expected: heredoc.Doc(`
				foo     bar
				foofoo  barbar
			`),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			output := New(out, err)

			output.Table(test.table)

			s.Equal(test.expected, out.String())
			s.Empty(err)
		})
	}
}
