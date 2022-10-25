package report

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReportSuite struct{ suite.Suite }

func TestReportSuiteSuite(t *testing.T) {
	suite.Run(t, new(ReportSuite))
}

func (s *ReportSuite) TestString() {

	s.Run("Empty", func() {
		report := &Report{}

		s.Empty(report.String())
	})

	s.Run("Message", func() {
		report := &Report{}
		report.Compose(
			WithMessage("message"),
		)

		s.Equal("message", report.String())
	})

	s.Run("Error", func() {
		report := &Report{}
		report.Compose(
			WithErr(fmt.Errorf("error")),
		)

		s.Equal("error", report.String())
	})

	s.Run("Message And Error", func() {
		report := &Report{}
		report.Compose(
			WithMessage("message"),
			WithErr(fmt.Errorf("error")),
		)

		s.Equal("message", report.String())
	})
}
