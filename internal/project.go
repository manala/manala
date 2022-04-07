package internal

import (
	_ "embed"
	"errors"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
	"github.com/imdario/mergo"
	"io"
	internalOs "manala/internal/os"
	internalTemplate "manala/internal/template"
	internalValidator "manala/internal/validator"
	internalYaml "manala/internal/yaml"
	"os"
	"path/filepath"
)

type Project struct {
	manifest *ProjectManifest
	recipe   *Recipe
}

func (project *Project) Path() string {
	return filepath.Dir(project.manifest.path)
}

func (project *Project) Manifest() *ProjectManifest {
	return project.manifest
}

func (project *Project) Recipe() *Recipe {
	return project.recipe
}

func (project *Project) Vars() map[string]interface{} {
	var vars map[string]interface{}

	_ = mergo.Merge(&vars, project.recipe.Vars())
	_ = mergo.Merge(&vars, project.manifest.Vars, mergo.WithOverride)

	return vars
}

func (project *Project) Template() *internalTemplate.Template {
	return project.recipe.Template().
		WithData(project)
}

func (project *Project) ManifestTemplate() *internalTemplate.Template {
	return project.recipe.ProjectManifestTemplate().
		WithData(project).
		WithDefaultContent(projectManifestTemplate)

}

/************/
/* Manifest */
/************/

//go:embed project_manifest.schema.json
var projectManifestSchema string

//go:embed project_manifest.yaml.tmpl
var projectManifestTemplate string

const projectDirLoaderManifest = ".manala.yaml"

func NewProjectManifest(dir string) *ProjectManifest {
	return &ProjectManifest{
		path: filepath.Join(dir, projectDirLoaderManifest),
		Vars: map[string]interface{}{},
	}
}

type ProjectManifest struct {
	path       string
	content    []byte
	Recipe     string
	Repository string
	Vars       map[string]interface{}
}

func (manifest *ProjectManifest) Path() string {
	return manifest.path
}

// Write implement the Writer interface
func (manifest *ProjectManifest) Write(content []byte) (n int, err error) {
	manifest.content = append(manifest.content, content...)
	return len(content), nil
}

// Load parse and decode content
func (manifest *ProjectManifest) Load() error {
	// Parse content
	contentFile, err := parser.ParseBytes(manifest.content, 0)
	if err != nil {
		return internalYaml.Error(manifest.path, err)
	}

	// Decode content
	var contentData map[string]interface{}
	if err := yaml.NewDecoder(contentFile).Decode(&contentData); err != nil {
		if err == io.EOF {
			return EmptyProjectManifestPathError(manifest.path)
		}
		return internalYaml.Error(manifest.path, err)
	}

	// Validate content
	if err, errs, ok := internalValidator.Validate(projectManifestSchema, contentData, internalValidator.WithYamlContent(manifest.content)); !ok {
		return ValidationProjectManifestPathError(manifest.path, err, errs)
	}

	// Decode manifest with "manala" data part
	manalaPath, _ := yaml.PathString("$.manala")
	manalaNode, _ := manalaPath.FilterFile(contentFile)
	if err := yaml.NewDecoder(manalaNode).Decode(manifest); err != nil {
		return internalYaml.Error(manifest.path, err)
	}

	// Keep remaining data for manifest vars
	delete(contentData, "manala")
	manifest.Vars = contentData

	return nil
}

func (manifest *ProjectManifest) Save() error {
	// Ensure manifest directory path exists
	path := filepath.Dir(manifest.path)
	if pathStat, err := os.Stat(manifest.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return internalOs.FileSystemError(err)
			}
		} else {
			return internalOs.FileSystemError(err)
		}
	} else if !pathStat.IsDir() {
		return WrongProjectManifestPathError(manifest.path)
	}

	if err := os.WriteFile(manifest.path, manifest.content, 0644); err != nil {
		return internalOs.FileSystemError(err)
	}

	return nil
}
