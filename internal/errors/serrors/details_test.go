package serrors

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DetailsSuite struct{ suite.Suite }

func TestDetailsSuite(t *testing.T) {
	suite.Run(t, new(DetailsSuite))
}

type DetailsError struct {
	*Details
}

func (s *DetailsSuite) Test() {
	err := &DetailsError{
		Details: &Details{},
	}

	s.Empty(err.ErrorDetails(false))

	err.SetDetails("details")

	s.Equal("details", err.ErrorDetails(false))
}
