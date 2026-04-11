package scale_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/image/scale"

	"github.com/stretchr/testify/suite"
)

type ScaleSuite struct {
	suite.Suite

	reference *image.NRGBA
}

func TestScaleSuite(t *testing.T) {
	suite.Run(t, new(ScaleSuite))
}

func (s *ScaleSuite) SetupTest() {
	s.reference = s.load(filepath.Join("testdata", "reference.png"))
}

func (s *ScaleSuite) Test() {
	tests := []struct {
		test     string
		expected *image.NRGBA
	}{
		{
			test:     "Same",
			expected: s.load(filepath.Join("testdata", "expected.8x8.png")),
		},
		{
			test:     "Bigger",
			expected: s.load(filepath.Join("testdata", "expected.12x12.png")),
		},
		{
			test:     "Smaller",
			expected: s.load(filepath.Join("testdata", "expected.4x4.png")),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			bounds := test.expected.Bounds()

			actual := scale.Scale(s.reference, bounds.Dx(), bounds.Dy())

			s.Require().Equal(bounds, actual.Bounds())

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					s.Equal(test.expected.NRGBAAt(x, y), actual.NRGBAAt(x, y))
				}
			}
		})
	}
}

func (s *ScaleSuite) load(path string) *image.NRGBA {
	s.T().Helper()

	file, _ := os.Open(path)
	img, _ := png.Decode(file)

	return img.(*image.NRGBA)
}
