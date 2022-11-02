package recipe

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"manala/core"
	"testing"
)

type NameProcessorManagerSuite struct{ suite.Suite }

func TestNameProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(NameProcessorManagerSuite))
}

func (s *NameProcessorManagerSuite) TestProcessName() {
	tests := []struct {
		name          string
		uppermostName string
		expected      string
		error         bool
	}{
		{
			name:          "",
			uppermostName: "",
			error:         true,
		},
		{
			name:          "name",
			uppermostName: "",
			expected:      "name",
		},
		{
			name:          "",
			uppermostName: "upper",
			expected:      "upper",
		},
		{
			name:          "name",
			uppermostName: "upper",
			expected:      "upper",
		},
	}

	for i, test := range tests {
		s.Run(fmt.Sprint(i), func() {
			manager := &NameProcessorManager{
				uppermostName: test.uppermostName,
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
