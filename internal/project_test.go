package internal

import (
	"bytes"
	_ "embed"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type ProjectSuite struct{ suite.Suite }

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

var projectTestPath = filepath.Join("testdata", "project")

func (s *ProjectSuite) Test() {
	repository := &Repository{path: "repository"}
	recipeManifest := NewRecipeManifest("dir")
	recipe := &Recipe{
		name:       "recipe",
		manifest:   recipeManifest,
		repository: repository,
	}
	projectManifest := NewProjectManifest("dir")
	project := &Project{
		manifest: projectManifest,
		recipe:   recipe,
	}

	s.Equal("dir", project.Path())
	s.Equal(projectManifest, project.Manifest())
	s.Equal(recipe, project.Recipe())

	s.Run("Vars", func() {
		recipeManifest.Vars = map[string]interface{}{
			"foo": "recipe",
			"bar": "recipe",
		}
		projectManifest.Vars = map[string]interface{}{
			"bar": "project",
			"baz": "project",
		}

		s.Equal(map[string]interface{}{
			"foo": "recipe",
			"bar": "project",
			"baz": "project",
		}, project.Vars())
	})

	s.Run("Template", func() {
		template := project.Template()

		out := &bytes.Buffer{}
		err := template.
			WithDefaultContent(`{{ .Vars | toYaml }}`).
			Write(out)

		s.NoError(err)
		s.Equal(`bar: project
baz: project
foo: recipe`, out.String())
	})

	s.Run("ManifestTemplate", func() {
		template := project.ManifestTemplate()

		out := &bytes.Buffer{}
		err := template.
			Write(out)

		s.NoError(err)
		s.Equal(`####################################################################
#                         !!! REMINDER !!!                         #
# Don't forget to run `+"`manala up`"+` each time you update this file ! #
####################################################################

manala:
    recipe: recipe
    repository: repository

# Default vars
bar: project
baz: project
foo: recipe
`, out.String())
	})
}

func (s *ProjectSuite) TestManifest() {
	projectManifest := NewProjectManifest("dir")

	s.Equal(filepath.Join("dir", ".manala.yaml"), projectManifest.Path())

	s.Run("Write", func() {
		length, err := projectManifest.Write([]byte("foo"))
		s.NoError(err)
		s.Equal(3, length)
		s.Equal("foo", string(projectManifest.content))
	})
}

func (s *ProjectSuite) TestManifestLoad() {
	projectManifest := NewProjectManifest("")

	path := filepath.Join(projectTestPath, "manifest_load")

	file, _ := os.Open(filepath.Join(path, "manifest.yaml"))
	_, _ = io.Copy(projectManifest, file)
	err := projectManifest.Load()

	s.NoError(err)
	s.Equal("recipe", projectManifest.Recipe)
	s.Equal("repository", projectManifest.Repository)
	s.Equal(map[string]interface{}{"foo": "bar"}, projectManifest.Vars)

	s.Run("Invalid Yaml", func() {
		projectManifest := NewProjectManifest("")

		file, _ := os.Open(filepath.Join(path, "manifest_invalid_yaml.yaml"))
		_, _ = io.Copy(projectManifest, file)
		err := projectManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
	})

	s.Run("Empty", func() {
		projectManifest := NewProjectManifest("")

		file, _ := os.Open(filepath.Join(path, "manifest_empty.yaml"))
		_, _ = io.Copy(projectManifest, file)
		err := projectManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("empty project manifest", internalError.Message)
	})

	s.Run("Wrong", func() {
		projectManifest := NewProjectManifest("")

		file, _ := os.Open(filepath.Join(path, "manifest_wrong.yaml"))
		_, _ = io.Copy(projectManifest, file)
		err := projectManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
	})

	s.Run("Invalid", func() {
		projectManifest := NewProjectManifest("")

		file, _ := os.Open(filepath.Join(path, "manifest_invalid.yaml"))
		_, _ = io.Copy(projectManifest, file)
		err := projectManifest.Load()

		s.ErrorAs(err, &internalError)
		s.Equal("project validation error", internalError.Message)
	})
}

func (s *ProjectSuite) TestManifestSave() {
	path := filepath.Join(projectTestPath, "manifest_save")

	_ = os.Remove(filepath.Join(path, ".manala.yaml"))
	_ = os.RemoveAll(filepath.Join(path, "directory"))

	projectManifest := NewProjectManifest(path)

	file, _ := os.Open(filepath.Join(path, "manifest.yaml"))
	_, _ = io.Copy(projectManifest, file)

	err := projectManifest.Save()

	s.NoError(err)
	s.FileExists(filepath.Join(path, ".manala.yaml"))
	sourceContent, _ := os.ReadFile(filepath.Join(path, "manifest.yaml"))
	destinationContent, _ := os.ReadFile(filepath.Join(path, ".manala.yaml"))
	s.Equal(sourceContent, destinationContent)

	s.Run("Directory", func() {
		directoryPath := filepath.Join(path, "directory")

		projectManifest := NewProjectManifest(directoryPath)

		file, _ := os.Open(filepath.Join(path, "manifest.yaml"))
		_, _ = io.Copy(projectManifest, file)

		err := projectManifest.Save()

		s.NoError(err)
		s.DirExists(directoryPath)
		s.FileExists(filepath.Join(directoryPath, ".manala.yaml"))
		sourceContent, _ := os.ReadFile(filepath.Join(path, "manifest.yaml"))
		destinationContent, _ := os.ReadFile(filepath.Join(directoryPath, ".manala.yaml"))
		s.Equal(sourceContent, destinationContent)
	})
}
