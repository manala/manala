package engine_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/template/engine"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type ExecutorSuite struct {
	suite.Suite

	engine *engine.Engine
	buffer *bytes.Buffer
}

func TestExecutorSuite(t *testing.T) {
	suite.Run(t, new(ExecutorSuite))
}

func (s *ExecutorSuite) SetupSuite() {
	s.engine = engine.New()
}

func (s *ExecutorSuite) SetupTest() {
	s.buffer = &bytes.Buffer{}
}

func (s *ExecutorSuite) TestExecute() {
	s.Run("Empty", func() {
		s.buffer.Reset()

		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, "")

		s.Require().NoError(err)
		s.Empty(s.buffer.String())
	})

	s.Run("Content", func() {
		s.buffer.Reset()

		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, `{{ "baz" }}`)

		s.Require().NoError(err)
		s.Equal("baz", s.buffer.String())
	})

	s.Run("Data", func() {
		s.buffer.Reset()

		executor, err := s.engine.Executor("foo")
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, `{{ . }}`)

		s.Require().NoError(err)
		s.Equal("foo", s.buffer.String())
	})

	s.Run("Partials", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecute/Partials")
		executor, err := s.engine.Executor(nil, filepath.Join(dir, "partial.tmpl"))
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, `{{ template "foo" . }}`)

		s.Require().NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("Include", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecute/Include")
		executor, err := s.engine.Executor(nil, filepath.Join(dir, "partial.tmpl"))
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, `{{ include "foo" . }}`)

		s.Require().NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("ParsingError", func() {
		s.buffer.Reset()

		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, "foo\n  {{ .bar }\nbaz\n")

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse template",
			Dump: heredoc.Doc(`
				  1 │ foo
				▶ 2 │   {{ .bar }
				    ├ unexpected "}" in operand
				  3 │ baz
			`),
		}, err)
	})

	s.Run("ExecutionError", func() {
		s.buffer.Reset()

		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.Execute(s.buffer, "foo\n  {{ .bar }}\nbaz\n")

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse template",
			Dump: heredoc.Doc(`
				  1 │ foo
				▶ 2 │   {{ .bar }}
				    ├──────╯ nil data; no entry for key "bar"
				  3 │ baz
			`),
		}, err)
	})
}

func (s *ExecutorSuite) TestExecuteTemplate() {
	s.Run("Empty", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/Empty")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Empty(s.buffer.String())
	})

	s.Run("Content", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/Content")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Equal("baz\n", s.buffer.String())
	})

	s.Run("Data", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/Data")
		executor, err := s.engine.Executor("foo")
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Equal("foo\n", s.buffer.String())
	})

	s.Run("Partials", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/Partials")
		executor, err := s.engine.Executor(nil, filepath.Join(dir, "partial.tmpl"))
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Equal("bar\n", s.buffer.String())
	})

	s.Run("Include", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/Include")
		executor, err := s.engine.Executor(nil, filepath.Join(dir, "partial.tmpl"))
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Equal("bar\n", s.buffer.String())
	})

	s.Run("File", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/File")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		s.Require().NoError(err)
		s.Equal("bar\n", s.buffer.String())
	})

	s.Run("FileNotFound", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/FileNotFound")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "not_found.tmpl"))

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to read template file",
			Attrs: [][2]any{
				{"file", filepath.Join(dir, "not_found.tmpl")},
			},
			Errors: []expect.ErrorExpectation{
				serrors.Expectation{
					Message: "file does not exist",
					Attrs: [][2]any{
						{"operation", "open"},
						{"path", filepath.Join(dir, "not_found.tmpl")},
					},
				},
			},
		}, err)
	})

	s.Run("FileInvalid", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/FileInvalid")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse template file",
			Dump: heredoc.Doc(`
				at %[1]s:2

				  1 │ foo
				▶ 2 │   {{ .bar }
				    ├ unexpected "}" in operand
				  3 │ baz
			`,
				filepath.Join(dir, "template.tmpl"),
			),
		}, err)
	})

	s.Run("ExecutionError", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/ExecutorSuite/TestExecuteTemplate/ExecutionError")
		executor, err := s.engine.Executor(nil)
		s.Require().NoError(err)

		err = executor.ExecuteTemplate(s.buffer, filepath.Join(dir, "template.tmpl"))

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse template file",
			Dump: heredoc.Doc(`
				at %[1]s:2:6

				  1 │ foo
				▶ 2 │   {{ .bar }}
				    ├──────╯ nil data; no entry for key "bar"
				  3 │ baz
			`,
				filepath.Join(dir, "template.tmpl"),
			),
		}, err)
	})
}
