package asciify_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/image/asciify"

	"github.com/stretchr/testify/suite"
)

type AsciifySuite struct {
	suite.Suite
}

func TestAsciifySuite(t *testing.T) {
	suite.Run(t, new(AsciifySuite))
}

func (s *AsciifySuite) Test() {
	tests := []struct {
		test      string
		reference *image.NRGBA
		expected  string
	}{
		{
			test:      "Opaque",
			reference: s.load(filepath.Join("testdata", "opaque.png")),
			expected: "" +
				"\x1b[38;2;255;0;0;48;2;0;0;255m▀▀" +
				"\x1b[38;2;0;255;0;48;2;255;255;255m▀▀" +
				"\n\x1b[m",
		},
		{
			test:      "Mixed",
			reference: s.load(filepath.Join("testdata", "mixed.png")),
			expected: "" +
				"\x1b[38;2;255;0;0;48;2;0;0;255m▀" +
				"\x1b[m " +
				"\x1b[38;2;255;0;0;49m▄" +
				"\x1b[38;2;0;255;0m▀" +
				"\n\x1b[m",
		},
		{
			test:      "Transparent",
			reference: s.load(filepath.Join("testdata", "transparent.png")),
			expected:  "    \n",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			actual := asciify.Asciify(test.reference)

			s.Equal(test.expected, actual)
		})
	}
}

func (s *AsciifySuite) load(path string) *image.NRGBA {
	s.T().Helper()

	file, _ := os.Open(path)
	img, _ := png.Decode(file)

	return img.(*image.NRGBA)
}
