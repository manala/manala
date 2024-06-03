package charm_test

import (
	"io"
	"manala/internal/ui/adapters/charm"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/suite"
)

type StylesSuite struct{ suite.Suite }

func TestStylesSuite(t *testing.T) {
	suite.Run(t, new(StylesSuite))
}

func (s *StylesSuite) TestFit() {
	definition := charm.NewStyleDefinition(
		lipgloss.NewStyle(),
	)
	renderer := lipgloss.NewRenderer(io.Discard)

	tests := []struct {
		test     string
		str      string
		width    int
		height   int
		expected string
	}{
		{
			test:     "Width",
			str:      "foo",
			width:    6,
			height:   0,
			expected: "foo   ",
		},
		{
			test:     "Height",
			str:      "foo",
			width:    0,
			height:   3,
			expected: "foo\n   \n   ",
		},
		{
			test:     "WidthAndHeight",
			str:      "foo",
			width:    6,
			height:   3,
			expected: "foo   \n      \n      ",
		},
		{
			test:     "WidthCrop",
			str:      "foo",
			width:    1,
			height:   0,
			expected: "f",
		},
		{
			test:     "HeightCrop",
			str:      "foo\nbar\nbaz",
			width:    0,
			height:   1,
			expected: "foo",
		},
		{
			test:     "WidthAndHeightCrop",
			str:      "foo\nbar\nbaz",
			width:    1,
			height:   1,
			expected: "f",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			style := definition.New(renderer)

			s.Equal(
				test.expected,
				style.Fit(test.str, test.width, test.height),
			)
		})
	}
}
