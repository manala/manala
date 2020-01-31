package cleaner

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/*****************/
/* Clean - Suite */
/*****************/

type CleanTestSuite struct{ suite.Suite }

func TestCleanTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(CleanTestSuite))
}

/*****************/
/* Clean - Tests */
/*****************/

func (s *CleanTestSuite) TestClean() {
	testMap := map[string]interface{}{
		"foo": "bar",
		"bar": map[string]interface{}{
			"foo": "bar",
		},
		"baz": map[interface{}]interface{}{
			"foo": "bar",
		},
	}
	testMap = Clean(testMap)
	s.Equal(
		map[string]interface{}{
			"foo": "bar",
			"bar": map[string]interface{}{
				"foo": "bar",
			},
			"baz": map[string]interface{}{
				"foo": "bar",
			},
		},
		testMap,
	)
}
