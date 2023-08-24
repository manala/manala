package recipe

import (
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app/mocks"
	"manala/core"
	"manala/internal/errors/serrors"
	"testing"
)

type NameProcessorManagerSuite struct{ suite.Suite }

func TestNameProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(NameProcessorManagerSuite))
}

func (s *NameProcessorManagerSuite) TestProcessName() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingManagerMock := &mocks.RecipeManagerMock{}

	tests := []struct {
		test         string
		name         string
		names        map[int]string
		expectedName string
		expectedErr  *serrors.Assert
	}{
		{
			test: "1",
			name: "",
			names: map[int]string{
				10: "",
			},
			expectedErr: &serrors.Assert{
				Type:    &core.UnprocessableRecipeNameError{},
				Message: "unable to process recipe name",
			},
		},
		{
			test: "2",
			name: "name",
			names: map[int]string{
				10: "",
			},
			expectedName: "name",
		},
		{
			test: "3",
			name: "",
			names: map[int]string{
				10: "upper",
			},
			expectedName: "upper",
		},
		{
			test: "4",
			name: "name",
			names: map[int]string{
				10: "upper",
			},
			expectedName: "upper",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			manager := NewNameProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, name := range test.names {
				manager.AddName(name, priority)
			}

			actual, err := manager.processName(test.name)

			if test.expectedErr != nil {
				s.Empty(actual)
				serrors.Equal(s.Assert(), test.expectedErr, err)
			} else {
				s.Equal(test.expectedName, actual)
				s.NoError(err)
			}
		})
	}
}
