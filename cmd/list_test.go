package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

/****************/
/* List - Suite */
/****************/

type ListTestSuite struct {
	suite.Suite
}

func TestListTestSuite(t *testing.T) {
	// Config
	viper.SetDefault("repository", "testdata/list/repository/default")
	// Run
	suite.Run(t, new(ListTestSuite))
}

/****************/
/* List - Tests */
/****************/

func (s *ListTestSuite) Test() {
	for _, t := range []struct {
		test   string
		args   []string
		err    string
		stdErr string
		stdOut string
	}{
		{
			test: "Default repository",
			args: []string{},
			err:  "",
			stdOut: `bar: Default bar recipe
foo: Default foo recipe
`,
			stdErr: "",
		},
		{
			test: "Use repository",
			args: []string{"--repository", "testdata/list/repository/custom"},
			err:  "",
			stdOut: `bar: Custom bar recipe
foo: Custom foo recipe
`,
			stdErr: "",
		},
		{
			test: "Use invalid repository",
			args: []string{"--repository", "testdata/list/repository/invalid"},
			err:  "\"testdata/list/repository/invalid\" directory does not exists",
		},
	} {
		s.Run(t.test, func() {
			// Command
			cmd := ListCmd()

			// Io
			stdOut := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			stdErr := bytes.NewBufferString("")
			cmd.SetErr(stdErr)
			log.SetHandler(cli.New(stdErr))

			// Execute
			cmd.SetArgs(t.args)
			err := cmd.Execute()

			// Test error
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)

				// Test stdout
				if t.stdOut == "" {
					s.Zero(stdOut.Len())
				} else {
					s.Equal(t.stdOut, stdOut.String())
				}

				// Test stderr
				if t.stdErr == "" {
					s.Zero(stdErr.Len())
				} else {
					s.Equal(t.stdErr, stdErr.String())
				}
			}
		})
	}
}
