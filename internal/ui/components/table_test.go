package components

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TableSuite struct {
	suite.Suite
}

func TestTableSuite(t *testing.T) {
	suite.Run(t, new(TableSuite))
}

func (s *TableSuite) Test() {
	table := &Table{}

	table.AddRow("primary 1", "secondary 1")
	table.AddRow("primary 2", "secondary 2")

	s.Equal([]*TableRow{
		{Primary: "primary 1", Secondary: "secondary 1"},
		{Primary: "primary 2", Secondary: "secondary 2"},
	}, table.Rows)
}
