package recipe

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/app/mocks"
	"manala/core"
	internalLog "manala/internal/log"
	"testing"
)

type NameProcessorManagerSuite struct{ suite.Suite }

func TestNameProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(NameProcessorManagerSuite))
}

func (s *NameProcessorManagerSuite) TestProcessName() {
	log := internalLog.New(io.Discard)

	cascadingManagerMock := mocks.MockRecipeManager()

	tests := []struct {
		name     string
		names    map[int]string
		expected string
		error    bool
	}{
		{
			name: "",
			names: map[int]string{
				10: "",
			},
			error: true,
		},
		{
			name: "name",
			names: map[int]string{
				10: "",
			},
			expected: "name",
		},
		{
			name: "",
			names: map[int]string{
				10: "upper",
			},
			expected: "upper",
		},
		{
			name: "name",
			names: map[int]string{
				10: "upper",
			},
			expected: "upper",
		},
	}

	for i, test := range tests {
		s.Run(fmt.Sprint(i), func() {
			manager := NewNameProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, name := range test.names {
				manager.AddName(name, priority)
			}

			actual, err := manager.processName(test.name)

			if test.error {
				var _error *core.UnprocessableRecipeNameError
				s.ErrorAs(err, &_error)
				s.Empty(actual)
			} else {
				s.NoError(err)
				s.Equal(test.expected, actual)
			}
		})
	}
}
