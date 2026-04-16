package name_test

import (
	"io"
	"testing"

	"github.com/manala/manala/app/recipe/name"
	"github.com/manala/manala/internal/log"

	"github.com/stretchr/testify/suite"
)

type ProcessorSuite struct{ suite.Suite }

func TestProcessorSuite(t *testing.T) {
	suite.Run(t, new(ProcessorSuite))
}

func (s *ProcessorSuite) TestProcess() {
	tests := []struct {
		test     string
		name     string
		names    map[int]string
		expected string
	}{
		{
			test: "1",
			name: "",
			names: map[int]string{
				10: "",
			},
		},
		{
			test: "2",
			name: "name",
			names: map[int]string{
				10: "",
			},
			expected: "name",
		},
		{
			test: "3",
			name: "",
			names: map[int]string{
				10: "upper",
			},
			expected: "upper",
		},
		{
			test: "4",
			name: "name",
			names: map[int]string{
				10: "upper",
			},
			expected: "upper",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			processor := name.NewProcessor(log.New(io.Discard))

			for weight, name := range test.names {
				processor.Add(name, weight)
			}

			name := processor.Process(test.name)

			s.Equal(test.expected, name)
		})
	}
}
