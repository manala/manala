package template_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/template"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type EngineSuite struct{ suite.Suite }

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}

func (s *EngineSuite) TestExecutor() {
	engine := template.NewEngine()

	repositoryMock := &app.RepositoryMock{}
	repositoryMock.
		On("URL").Return("url")

	recipeMock := &app.RecipeMock{}
	recipeMock.
		On("Name").Return("name").
		On("Description").Return("description").
		On("Icon").Return("icon").
		On("Repository").Return(repositoryMock).
		On("Partials").Return([]string{})

	executor, err := engine.Executor(map[string]any{"foo": "bar"}, recipeMock, "dir")
	s.Require().NoError(err)

	buffer := &bytes.Buffer{}
	err = executor.Execute(buffer, strings.TrimLeft(`
.Dir: {{ .Dir }}
.Vars: {{ .Vars | toJson }}
.Recipe.Name: {{ .Recipe.Name }}
.Recipe.Description: {{ .Recipe.Description }}
.Recipe.Icon: {{ .Recipe.Icon }}
.Recipe.Repository.URL: {{ .Recipe.Repository.URL }}
.Recipe.Repository.Path: {{ .Recipe.Repository.Path }}
.Recipe.Repository.Source: {{ .Recipe.Repository.Source }}
.Repository.URL: {{ .Repository.URL }}
`, "\n"))
	s.Require().NoError(err)

	heredoc.Equal(s.T(), `
		.Dir: dir
		.Vars: {"foo":"bar"}
		.Recipe.Name: name
		.Recipe.Description: description
		.Recipe.Icon: icon
		.Recipe.Repository.URL: url
		.Recipe.Repository.Path: url
		.Recipe.Repository.Source: url
		.Repository.URL: url
	`, buffer.String())
}
